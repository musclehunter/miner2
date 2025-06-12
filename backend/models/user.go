package models

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User はユーザー情報を表す構造体
type User struct {
	ID                    string    `json:"id" gorm:"primaryKey"`
	Email                 string    `json:"email" gorm:"unique;not null"`
	PasswordHash          string    `json:"-" gorm:"column:password_hash;not null"` // JSONレスポンスには含めない
	Salt                  string    `json:"-" gorm:"column:salt;not null"`          // JSONレスポンスには含めない
	Name                  string    `json:"name" gorm:"not null"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// generateSalt は新しいランダムなソルトを生成します
func generateSalt() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// SetPassword はパスワードをソルト付きでハッシュ化して設定
func (u *User) SetPassword(password string) error {
	// 新しいソルトを生成
	salt, err := generateSalt()
	if err != nil {
		return err
	}
	u.Salt = salt
	
	// ソルト付きパスワードをハッシュ化
	saltedPassword := password + u.Salt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	
	u.PasswordHash = string(hashedPassword)
	return nil
}

// CheckPassword はパスワードが正しいかチェック
func (u *User) CheckPassword(password string) bool {
	// ソルト付きパスワードを検証
	saltedPassword := password + u.Salt
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(saltedPassword))
	return err == nil
}

// GenerateVerificationToken はメール確認用のトークンを生成します
// トークンと有効期限を返しますが、User構造体には保存しません
func GenerateVerificationToken() (string, time.Time, error) {
	// ランダムなトークンを生成
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", time.Time{}, err
	}
	
	// Base64でエンコード
	token := base64.URLEncoding.EncodeToString(bytes)
	
	// 有効期限を24時間後に設定
	expiry := time.Now().Add(24 * time.Hour)
	
	return token, expiry, nil
}

// VerifyEmailToken はトークンが有効期限内かチェックします
// この関数は別途保存されたトークンと有効期限を使用します
func VerifyEmailToken(token string, expiry time.Time) bool {
	// 有効期限内かチェック
	return time.Now().Before(expiry)
}

// NewUser は新しいユーザーを作成
func NewUser(email, name, password string) (*User, error) {
	user := &User{
		Email:     email,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}
	
	return user, nil
}
