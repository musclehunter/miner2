package models

import (
	"time"
	"golang.org/x/crypto/bcrypt"
)

// User はユーザー情報を表す構造体
type User struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"` // JSONレスポンスには含めない
	Name      string    `json:"name" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SetPassword はパスワードをハッシュ化して設定
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword はパスワードが正しいかチェック
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
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
