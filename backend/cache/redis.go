package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/musclehunter/miner2/models"
)

// キーのプレフィックス
const (
	EmailVerificationPrefix = "email_verification:"
)

// RedisClient はRedisとのインタラクションを管理するインターフェース
type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
	Exists(ctx context.Context, key string) (bool, error)
	Keys(ctx context.Context, pattern string) ([]string, error)
}

// MockRedisClient は開発環境用のモックRedisクライアント
type MockRedisClient struct {
	data map[string]mockRedisItem
}

// mockRedisItem はキャッシュアイテムを表す
type mockRedisItem struct {
	value      string
	expiration time.Time
}

// NewMockRedisClient は新しいモックRedisクライアントを作成
func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		data: make(map[string]mockRedisItem),
	}
}

// Set は値をキャッシュに設定
func (c *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	expireAt := time.Now().Add(expiration)
	c.data[key] = mockRedisItem{
		value:      fmt.Sprintf("%v", value),
		expiration: expireAt,
	}
	log.Printf("MockRedis: SET %s=%v (expires: %v)", key, value, expireAt)
	return nil
}

// Get はキャッシュから値を取得
func (c *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	item, exists := c.data[key]
	if !exists {
		return "", errors.New("key not found")
	}

	// 期限切れをチェック
	if time.Now().After(item.expiration) {
		delete(c.data, key)
		return "", errors.New("key expired")
	}

	log.Printf("MockRedis: GET %s=%s", key, item.value)
	return item.value, nil
}

// Del はキャッシュから値を削除
func (c *MockRedisClient) Del(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		delete(c.data, key)
		log.Printf("MockRedis: DEL %s", key)
	}
	return nil
}

// Exists はキーが存在するかチェック
func (c *MockRedisClient) Exists(ctx context.Context, key string) (bool, error) {
	item, exists := c.data[key]
	if !exists {
		return false, nil
	}

	// 期限切れをチェック
	if time.Now().After(item.expiration) {
		delete(c.data, key)
		return false, nil
	}

	return true, nil
}

// Keys はパターンに一致するキーを取得
func (c *MockRedisClient) Keys(ctx context.Context, pattern string) ([]string, error) {
	// この実装ではパターンマッチングは単純化 (prefix*)
	prefix := pattern[:len(pattern)-1] // 最後の「*」を削除
	var result []string

	for key, item := range c.data {
		// 期限切れをチェック
		if time.Now().After(item.expiration) {
			delete(c.data, key)
			continue
		}

		// プレフィックスマッチング
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			result = append(result, key)
		}
	}

	return result, nil
}

// Client はデフォルトのRedisクライアント
var Client RedisClient

// 実際のRedisクライアント実装は別ファイルで定義

// InitRedisClient はRedisクライアントを初期化
func InitRedisClient() {
	// 開発環境ではモッククライアントを使用
	// 本番環境ではgo-redis/redisクライアントを使用予定
	if os.Getenv("APP_ENV") == "production" {
		log.Println("本番環境用Redisクライアントを初期化します（未実装）")
		// 実際のRedisクライアントの初期化はまだ実装されていません
		// 現時点ではモッククライアントを使用します
		Client = NewMockRedisClient()
	} else {
		log.Println("開発環境用モックRedisクライアントを初期化します")
		Client = NewMockRedisClient()
	}
}

// メール確認関連の機能

// SaveEmailVerification はメール確認トークンを保存
func SaveEmailVerification(ctx context.Context, token, userInfo string) error {
	key := EmailVerificationPrefix + token
	// 24時間の有効期限
	return Client.Set(ctx, key, userInfo, 24*time.Hour)
}

// GetEmailVerification はメール確認トークンからユーザー情報を取得
func GetEmailVerification(ctx context.Context, token string) (string, error) {
	key := EmailVerificationPrefix + token
	return Client.Get(ctx, key)
}

// DeleteEmailVerification はメール確認トークンを削除
func DeleteEmailVerification(ctx context.Context, token string) error {
	key := EmailVerificationPrefix + token
	return Client.Del(ctx, key)
}

// CheckEmailExists はメールアドレスが仮登録に存在するか確認
func CheckEmailExists(ctx context.Context, email string) (bool, error) {
	// メール確認トークンに関連するすべてのキーを取得
	keys, err := Client.Keys(ctx, EmailVerificationPrefix+"*")
	if err != nil {
		return false, err
	}

	// 各キーに対応するユーザー情報を取得し、メールアドレスをチェック
	for _, key := range keys {
		value, err := Client.Get(ctx, key)
		if err != nil {
			log.Printf("キー %s の取得エラー: %v", key, err)
			continue
		}

		// 仮登録ユーザー情報を解析
		var pendingUser models.PendingUser
		if err := json.Unmarshal([]byte(value), &pendingUser); err != nil {
			log.Printf("ユーザー情報の解析エラー: %v", err)
			continue
		}

		// メールアドレスが一致するかチェック
		if pendingUser.Email == email {
			return true, nil
		}
	}

	return false, nil
}

// GetAllPendingUsers は全ての仮登録ユーザー情報を取得
func GetAllPendingUsers(ctx context.Context) ([]map[string]interface{}, error) {
	// メール確認トークンに関連するすべてのキーを取得
	keys, err := Client.Keys(ctx, EmailVerificationPrefix+"*")
	if err != nil {
		return nil, err
	}

	// 結果を格納するスライス
	var pendingUsers []map[string]interface{}

	// 各キーに対応するユーザー情報を取得
	for _, key := range keys {
		// キーからトークンを抽出
		token := key[len(EmailVerificationPrefix):]
		
		value, err := Client.Get(ctx, key)
		if err != nil {
			log.Printf("キー %s の取得エラー: %v", key, err)
			continue
		}

		// 仮登録ユーザー情報を解析
		var pendingUser models.PendingUser
		if err := json.Unmarshal([]byte(value), &pendingUser); err != nil {
			log.Printf("ユーザー情報の解析エラー: %v", err)
			continue
		}

		// 管理画面表示用にデータを変換
		userInfo := map[string]interface{}{
			"token": token,
			"email": pendingUser.Email,
			"name":  pendingUser.Name,
		}

		pendingUsers = append(pendingUsers, userInfo)
	}

	return pendingUsers, nil
}
