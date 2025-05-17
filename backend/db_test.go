package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestDBConnection(t *testing.T) {
	// 環境変数のチェック
	t.Log("環境変数のチェック:")
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

	t.Logf("DB_HOST=%s", dbHost)
	t.Logf("DB_PORT=%s", dbPort)
	t.Logf("DB_USER=%s", dbUser)
	t.Logf("DB_NAME=%s", dbName)

	// 接続文字列を構築
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	t.Logf("接続文字列: %s", dsn)

	// 接続を試みる
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("データベース接続エラー: %v", err)
	}
	defer db.Close()

	// 接続テスト
	if err = db.Ping(); err != nil {
		t.Fatalf("データベースPingエラー: %v", err)
	}

	t.Log("データベース接続に成功しました")

	// 現在時刻をデータベースから取得
	var dbTime string
	err = db.QueryRow("SELECT NOW()").Scan(&dbTime)
	if err != nil {
		t.Fatalf("データベースクエリエラー: %v", err)
	}
	t.Logf("データベース時刻: %s", dbTime)

	// テーブル一覧を取得
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		t.Fatalf("テーブル一覧取得エラー: %v", err)
	}
	defer rows.Close()

	t.Log("データベーステーブル一覧:")
	var tableName string
	for rows.Next() {
		if err := rows.Scan(&tableName); err != nil {
			t.Fatalf("テーブル名スキャンエラー: %v", err)
		}
		t.Logf("- %s", tableName)
	}

	// usersテーブルの確認
	rows, err = db.Query("SHOW COLUMNS FROM users")
	if err != nil {
		t.Logf("usersテーブルが存在しないか、エラーが発生しました: %v", err)
	} else {
		t.Log("usersテーブルのカラム構造:")
		
		var field, fieldType, null, key string
		var defaultValue, extra interface{}
		for rows.Next() {
			if err := rows.Scan(&field, &fieldType, &null, &key, &defaultValue, &extra); err != nil {
				t.Fatalf("カラムデータのスキャンエラー: %v", err)
			}
			t.Logf("- %s (%s, %s, key=%s)", field, fieldType, null, key)
		}
		rows.Close()
	}

	// ユーザーテーブルの作成（存在しない場合）
	// 注意: 既存のテーブル構造に合わせています
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(36) PRIMARY KEY,
		email VARCHAR(255) NOT NULL UNIQUE,
		password_hash VARCHAR(255) NOT NULL,
		salt VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)
	`)
	if err != nil {
		t.Fatalf("ユーザーテーブル作成エラー: %v", err)
	}
	t.Log("ユーザーテーブルを作成または確認しました")

	// テストユーザーの作成または確認
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", "test@example.com").Scan(&count)
	if err != nil {
		t.Fatalf("ユーザー検索エラー: %v", err)
	}

	if count == 0 {
		_, err = db.Exec(
			"INSERT INTO users (id, email, password_hash, salt, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())",
			"test-1", "test@example.com", "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", "test-salt",
		)
		if err != nil {
			t.Fatalf("テストユーザー作成エラー: %v", err)
		}
		t.Log("テストユーザーを作成しました: test@example.com")
	} else {
		t.Log("テストユーザーは既に存在します: test@example.com")
	}

	// ユーザー一覧の表示
	rows, err = db.Query("SELECT id, email, password_hash, created_at FROM users LIMIT 5")
	if err != nil {
		t.Fatalf("ユーザー一覧取得エラー: %v", err)
	}
	defer rows.Close()

	t.Log("ユーザー一覧:")
	var id, email, passHash, createdAt string
	for rows.Next() {
		if err := rows.Scan(&id, &email, &passHash, &createdAt); err != nil {
			t.Fatalf("ユーザーデータスキャンエラー: %v", err)
		}
		t.Logf("- %s: %s (パスワードハッシュ: %s...) 作成日: %s", id, email, passHash[:10], createdAt)
	}
}
