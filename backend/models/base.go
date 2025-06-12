package models

import "time"

// Base はプレイヤーの拠点を表します。
// `bases`テーブルに対応します。
type Base struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	TownID    string    `json:"town_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
