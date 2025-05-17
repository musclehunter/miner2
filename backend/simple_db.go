package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// MySQL データベース接続
var db *sql.DB

// データベース接続の初期化
func initDB() error {
	// 環境変数から接続情報を取得
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// 接続文字列を構築
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// 接続を試みる
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("データベース接続エラー: %v", err)
	}

	// 接続テスト
	if err = db.Ping(); err != nil {
		return fmt.Errorf("データベースPingエラー: %v", err)
	}

	return nil
}

// データベースヘルスチェックハンドラー
func dbHealthHandler(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "エラー",
			"error":  "データベース接続が初期化されていません",
		})
		return
	}

	// Ping テスト
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

// ユーザーテーブルのセットアップ
func setupUserTable() error {
	// ユーザーテーブルの作成
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(36) PRIMARY KEY,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		name VARCHAR(100) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)
	`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("ユーザーテーブル作成エラー: %v", err)
	}

	// テストユーザーの存在確認
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", "test@example.com").Scan(&count)
	if err != nil {
		return fmt.Errorf("ユーザー検索エラー: %v", err)
	}

	// テストユーザーが存在しない場合は作成
	if count == 0 {
		// パスワードをハッシュ化する処理は本来必要ですが、簡略化のため平文で保存
		_, err = db.Exec(
			"INSERT INTO users (id, email, password, name, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())",
			"1", "test@example.com", "password123", "テストユーザー",
		)
		if err != nil {
			return fmt.Errorf("テストユーザー作成エラー: %v", err)
		}
		log.Println("テストユーザーを作成しました: test@example.com")
	} else {
		log.Println("テストユーザーは既に存在します: test@example.com")
	}

	return nil
}

// トランザクション処理でユーザーを作成するサンプル
func createUserSample(email, name, password string) error {
	// トランザクション開始
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("トランザクション開始エラー: %v", err)
	}

	// ユーザーの存在確認
	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("ユーザー検索エラー: %v", err)
	}

	// ユーザーが既に存在する場合
	if count > 0 {
		tx.Rollback()
		return fmt.Errorf("このメールアドレスは既に使用されています: %s", email)
	}

	// 新しいユーザーを作成
	_, err = tx.Exec(
		"INSERT INTO users (id, email, password, name, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())",
		fmt.Sprintf("%d", time.Now().UnixNano()), email, password, name,
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("ユーザー作成エラー: %v", err)
	}

	// トランザクションをコミット
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("トランザクションコミットエラー: %v", err)
	}

	log.Printf("ユーザーを作成しました: %s", email)
	return nil
}
