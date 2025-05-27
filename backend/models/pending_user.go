package models

import (
	"encoding/json"
	"time"
)

// PendingUser は確認待ちユーザー情報を表す構造体
type PendingUser struct {
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	Password    string    `json:"password"` // ハッシュ化前のパスワード
	CreatedAt   time.Time `json:"created_at"`
}

// NewPendingUser は新しい確認待ちユーザーを作成
func NewPendingUser(email, name, password string) *PendingUser {
	return &PendingUser{
		Email:     email,
		Name:      name,
		Password:  password,
		CreatedAt: time.Now(),
	}
}

// ToJSON はユーザー情報をJSON文字列に変換
func (p *PendingUser) ToJSON() (string, error) {
	bytes, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FromJSON はJSON文字列からユーザー情報を復元
func PendingUserFromJSON(data string) (*PendingUser, error) {
	var p PendingUser
	err := json.Unmarshal([]byte(data), &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
