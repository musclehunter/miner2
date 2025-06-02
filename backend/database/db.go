package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DB はグローバルなデータベース接続を保持します
var DB *sql.DB

// エラー定義
var (
	ErrDatabaseNotInitialized = errors.New("データベースが初期化されていません")
)

// InitDB はデータベース接続を初期化します
func InitDB() error {
	// 環境変数からDB接続情報を取得（環境変数がない場合はデフォルト値を使用）
	dbHost := getEnv("DB_HOST", "db")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "miner_user")
	dbPass := getEnv("DB_PASSWORD", "miner_password")
	dbName := getEnv("DB_NAME", "miner_db")

	// MySQL接続文字列
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// 接続を試みる（最大10回、1秒間隔）
	var err error
	for i := 0; i < 10; i++ {
		DB, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("データベース接続エラー: %v, 再試行中...", err)
			time.Sleep(time.Second)
			continue
		}

		// 接続確認
		if err = DB.Ping(); err != nil {
			log.Printf("データベースPingエラー: %v, 再試行中...", err)
			time.Sleep(time.Second)
			continue
		}

		log.Println("データベース接続に成功しました")
		return nil
	}

	return fmt.Errorf("データベース接続に失敗しました: %v", err)
}

// getEnv は環境変数を取得し、設定されていない場合はデフォルト値を返します
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// CloseDB はデータベース接続を閉じます
func CloseDB() {
	if DB != nil {
		DB.Close()
		DB = nil
		log.Println("データベース接続を閉じました")
	}
}

// ConnectDB は新しいデータベース接続を返します（マイグレーション用）
func ConnectDB() (*sql.DB, error) {
	// 環境変数からDB接続情報を取得（環境変数がない場合はデフォルト値を使用）
	dbHost := getEnv("DB_HOST", "db")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "miner_user")
	dbPass := getEnv("DB_PASSWORD", "miner_password")
	dbName := getEnv("DB_NAME", "miner_db")

	// MySQL接続文字列
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// 接続確認
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
