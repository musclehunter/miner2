package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	
	"github.com/musclehunter/miner2/database"
	"github.com/musclehunter/miner2/handlers"
)

// User はユーザー情報を表す構造体
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Salt         string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ユーザーリポジトリ
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetUserByEmail gets a user by email
func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	query := "SELECT id, email, password_hash, salt, created_at, updated_at FROM users WHERE email = ?"
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Salt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ユーザーが見つかりません: %s", email)
		}
		return nil, fmt.Errorf("ユーザー検索エラー: %v", err)
	}
	return user, nil
}

// GetAllUsers gets all users
func (r *UserRepository) GetAllUsers() ([]*User, error) {
	query := "SELECT id, email, password_hash, salt, created_at, updated_at FROM users"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ユーザー一覧取得エラー: %v", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.Salt,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ユーザーデータスキャンエラー: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}

// ヘルスチェックハンドラー
func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "サーバーは正常に動作しています",
	})
}

// データベースヘルスチェックハンドラー
func dbHealthCheckHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if db == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "エラー",
				"error":  "データベース接続が初期化されていません",
			})
			return
		}

		// Pingテスト
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "エラー",
				"error":  fmt.Sprintf("データベース接続エラー: %v", err),
			})
			return
		}

		// 現在時刻をデータベースから取得
		var dbTime time.Time
		err := db.QueryRow("SELECT NOW()").Scan(&dbTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "エラー",
				"error":  fmt.Sprintf("データベースクエリエラー: %v", err),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "OK",
			"message": fmt.Sprintf("データベース接続成功: 現在時刻 %s (UTC)", dbTime.UTC().Format(time.RFC3339)),
		})
	}
}

// ユーザー一覧ハンドラー
func getUsersHandler(userRepo *UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := userRepo.GetAllUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "エラー",
				"error":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"users":  users,
		})
	}
}

// メイン関数
func runSimpleServer() {
	// 環境変数から接続情報を取得
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" {
		dbHost = "db" // デフォルト値
	}
	if dbPort == "" {
		dbPort = "3306" // デフォルト値
	}
	if dbUser == "" {
		dbUser = "miner" // デフォルト値
	}
	if dbPass == "" {
		dbPass = "minerpassword" // デフォルト値
	}
	if dbName == "" {
		dbName = "minerdb" // デフォルト値
	}

	// 接続文字列を構築
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	log.Printf("接続文字列: %s", dsn)

	// データベース接続
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("データベース接続エラー: %v", err)
	}
	defer db.Close()

	// 接続テスト
	if err = db.Ping(); err != nil {
		log.Fatalf("データベースPingエラー: %v", err)
	}
	log.Println("データベース接続に成功しました")

	// データベースをグローバル変数に設定してhandlersから使用可能にする
	database.DB = db

	// 初期データの作成
	log.Println("初期データをチェックしています...")
	if err := database.CreateInitialTowns(); err != nil {
		log.Printf("警告: 初期町データ作成エラー: %v", err)
	}
	if err := database.CreateInitialOres(); err != nil {
		log.Printf("警告: 初期鉱石データ作成エラー: %v", err)
	}

	// ハンドラー初期化 - 重要: handlers内で共有されるUserRepositoryが初期化される
	handlers.InitHandlers()   // 認証ハンドラー初期化
	handlers.InitGameHandlers() // ゲームハンドラー初期化
	
	// ユーザーリポジトリの作成 - こちらは/api/usersエンドポイント専用
	userRepo := NewUserRepository(db)

	// Ginルーターの設定
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	// CORS設定
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// ヘルスチェックエンドポイント
	r.GET("/ping", healthCheckHandler)
	r.GET("/api/health", healthCheckHandler)
	r.GET("/api/health/db", dbHealthCheckHandler(db))

	// ユーザー一覧エンドポイント
	r.GET("/api/users", getUsersHandler(userRepo))
	
	// 認証関連のエンドポイント
	auth := r.Group("/api/auth")
	{
		auth.POST("/signup", handlers.Signup)
		auth.POST("/login", handlers.Login)
		auth.GET("/me", handlers.AuthMiddleware(), handlers.Me)
	}
	
	// ゲーム関連のエンドポイント
	game := r.Group("/api/game")
	{
		// 町情報
		game.GET("/towns", handlers.GetAllTowns)
		game.GET("/towns/:id", handlers.GetTownByID)
		
		// 鉱石情報
		game.GET("/ores", handlers.GetAllOres)
		game.GET("/ores/:id", handlers.GetOreByID)
		
		// 認証が必要なエンドポイント
		secured := game.Group("/")
		secured.Use(handlers.AuthMiddleware())
		{
			// 今後ここに認証が必要なゲームエンドポイントを追加
		}
	}

	// サーバー起動
	port := "8080"
	log.Printf("サーバーを起動しています: http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("サーバー起動エラー: %v", err)
	}
}
