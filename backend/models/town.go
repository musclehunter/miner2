package models

import (
	"time"
)

// Town は町の情報を表す構造体
type Town struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewTown は新しい町を作成
func NewTown(id string, name string, description string) *Town {
	now := time.Now()
	return &Town{
		ID:          id,
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
