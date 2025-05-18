package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"

	"github.com/musclehunter/miner2/database"
	"github.com/musclehunter/miner2/models"
)

// ユーザーリポジトリ
var userRepo *database.UserRepository

// InitHandlersはハンドラの初期化を行います
func InitHandlers() {
	userRepo = database.NewUserRepository(database.DB)
	log.Println("認証ハンドラーの初期化完了: UserRepository設定済み")
}

// JWTの秘密鍵（本番環境では環境変数から取得するべき）
var jwtSecret = []byte("your-secret-key")

// サインアップリクエスト
type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

// ログインリクエスト
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// JWTクレーム
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// Signup_CreateMockUser はテスト用のユーザーを作成（開発用）
func Signup_CreateMockUser(email, name, password string) (*models.User, error) {
	// 既に存在するかチェック
	existingUser, err := userRepo.GetUserByEmail(email)
	if err != nil {
		log.Printf("ユーザー検索エラー: %v", err)
		return nil, err
	}

	if existingUser != nil {
		// 既に存在する場合はそのユーザーを返す
		log.Printf("ユーザーが既に存在します: %s", email)
		return existingUser, nil
	}

	// 新しいユーザーを作成
	user, err := models.NewUser(email, name, password)
	if err != nil {
		log.Printf("ユーザー作成エラー: %v", err)
		return nil, err
	}

	// UUIDを生成
	user.ID = uuid.New().String()

	// データベースに保存
	err = userRepo.CreateUser(user)
	if err != nil {
		log.Printf("ユーザー保存エラー: %v", err)
		return nil, err
	}

	log.Printf("モックユーザーを作成しました: %s", email)

	return user, nil
}

// Signup はユーザー登録を処理
func Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力データが不正です", "details": err.Error()})
		return
	}

	// メールアドレスが既に使用されているかチェック
	existingUser, err := userRepo.GetUserByEmail(req.Email)
	if err != nil {
		log.Printf("ユーザー検索エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー検索中にエラーが発生しました"})
		return
	}

	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "このメールアドレスは既に使用されています"})
		return
	}

	// 新しいユーザーを作成
	user, err := models.NewUser(req.Email, req.Name, req.Password)
	if err != nil {
		log.Printf("ユーザー作成エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー作成に失敗しました"})
		return
	}

	// UUIDを生成
	user.ID = uuid.New().String()

	// ユーザーをデータベースに保存
	err = userRepo.CreateUser(user)
	if err != nil {
		log.Printf("ユーザー保存エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー保存に失敗しました"})
		return
	}

	// JWTトークンを生成
	token, err := generateToken(user.ID)
	if err != nil {
		log.Printf("トークン生成エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "トークン生成に失敗しました"})
		return
	}

	log.Printf("新規ユーザー登録成功: %s", req.Email)

	// レスポンス
	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

// Login はユーザーログインを処理
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力データが不正です", "details": err.Error()})
		return
	}

	log.Printf("ログイン試行: %s", req.Email)

	// ユーザーをデータベースから検索
	log.Printf("ユーザー検索開始: email=%s, userRepo=%v", req.Email, userRepo)
	foundUser, err := userRepo.GetUserByEmail(req.Email)
	if err != nil {
		log.Printf("ログイン時のユーザー検索エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ログイン処理中にエラーが発生しました"})
		return
	}

	// ユーザーが見つからない場合
	if foundUser == nil {
		log.Printf("ログイン失敗: ユーザーが見つかりません %s", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "メールアドレスまたはパスワードが正しくありません"})
		return
	}
	
	log.Printf("ユーザー検索結果: %+v", foundUser)

	// パスワードを検証
	if !foundUser.CheckPassword(req.Password) {
		log.Printf("ログイン失敗: パスワードが一致しません %s", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "メールアドレスまたはパスワードが正しくありません"})
		return
	}

	// JWTトークンを生成
	log.Printf("トークン生成開始: ユーザーID=%s", foundUser.ID)
	token, err := generateToken(foundUser.ID)
	if err != nil {
		log.Printf("トークン生成エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "トークン生成に失敗しました"})
		return
	}

	log.Printf("ログイン成功: メール=%s, ユーザーID=%s, トークン=%s", req.Email, foundUser.ID, token[:10]+"...")

	// レスポンス
	response := gin.H{
		"token": token,
		"user": gin.H{
			"id":    foundUser.ID,
			"email": foundUser.Email,
			"name":  foundUser.Name,
		},
	}
	log.Printf("レスポンス送信: %+v", response)
	c.JSON(http.StatusOK, response)
}

// generateToken はJWTトークンを生成
func generateToken(userID string) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Printf("トークン生成エラー: %v", err)
		return "", err
	}

	return tokenString, nil
}

// AuthMiddleware は認証ミドルウェア
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証トークンがありません"})
			c.Abort()
			return
		}

		// "Bearer "の部分を削除
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// トークンを解析
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無効なトークンです"})
			c.Abort()
			return
		}

		// ユーザーIDをコンテキストに設定
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

// Me は現在のログインユーザー情報を取得
func Me(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証されていません"})
		return
	}

	// データベースからユーザー情報を取得
	user, err := userRepo.GetUserByID(userID.(string))
	if err != nil {
		log.Printf("ユーザー情報取得エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー情報の取得に失敗しました"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}
