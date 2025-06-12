package models

import "time"

// Warehouse は拠点の倉庫を表します。
// `warehouses`テーブルに対応します。
type Warehouse struct {
	ID        string    `json:"id"`
	BaseID    string    `json:"base_id"`
	Level     int       `json:"level"`
	Capacity  int       `json:"capacity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
