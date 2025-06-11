package models

import "time"

// PlayerBase はプレイヤーの拠点情報を表します。
type PlayerBase struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	TownID    string    `json:"town_id"`
	Level     int       `json:"level"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
