package models

import "time"

// PlayerInventory はプレイヤーの所持金と在庫容量を管理します。
type PlayerInventory struct {
	ID              string    `json:"id" db:"id"`
	UserID          string    `json:"user_id" db:"user_id"`
	Gold            int       `json:"gold" db:"gold"`
	MaxCapacity     int       `json:"max_capacity" db:"max_capacity"`
	CurrentCapacity int       `json:"current_capacity" db:"current_capacity"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	// User            User      `json:"user"` // JOIN時に設定
}
