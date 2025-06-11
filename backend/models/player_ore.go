package models

import "time"

// PlayerOre はプレイヤーが所持している鉱石の数量を管理します。
type PlayerOre struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	OreID     string    `json:"ore_id" db:"ore_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	// User      User      `json:"user"` // JOIN時に設定
	Ore       Ore       `json:"ore"` // JOIN時に設定
}
