package models

import "time"

// PlayerItem はプレイヤーが所持しているアイテムの数量を管理します。
type PlayerItem struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	ItemID    string    `json:"item_id" db:"item_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	// User      User      `json:"user"` // JOIN時に設定
	Item      Item      `json:"item"` // JOIN時に設定
}
