package models

import (
	"time"
)

// Ore は鉱石の情報を表す構造体
type Ore struct {
	ID                  string    `json:"id" gorm:"primaryKey"`
	Name                string    `json:"name" gorm:"not null"`
	Rarity              int       `json:"rarity" gorm:"not null"`
	Purity              float64   `json:"purity" gorm:"not null"`
	ProcessingDifficulty int       `json:"processing_difficulty" gorm:"not null"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// NewOre は新しい鉱石を作成
func NewOre(id string, name string, rarity int, purity float64, processingDifficulty int) *Ore {
	now := time.Now()
	return &Ore{
		ID:                  id,
		Name:                name,
		Rarity:              rarity,
		Purity:              purity,
		ProcessingDifficulty: processingDifficulty,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
}
