package database

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/musclehunter/miner2/models"
)

// PlayerInventoryRepository はプレイヤーの在庫に関連するデータベース操作を管理します
type PlayerInventoryRepository struct {
	db *sql.DB
}

// NewPlayerInventoryRepository は新しい PlayerInventoryRepository インスタンスを作成します
func NewPlayerInventoryRepository(db *sql.DB) *PlayerInventoryRepository {
	return &PlayerInventoryRepository{db: db}
}

// GetInventoryByUserID はユーザーIDで在庫情報を取得します
func (r *PlayerInventoryRepository) GetInventoryByUserID(userID string) (*models.PlayerInventory, error) {
	query := `
		SELECT id, user_id, gold, max_capacity, current_capacity, created_at, updated_at
		FROM player_inventories
		WHERE user_id = ?
	`

	var inventory models.PlayerInventory
	err := r.db.QueryRow(query, userID).Scan(
		&inventory.ID,
		&inventory.UserID,
		&inventory.Gold,
		&inventory.MaxCapacity,
		&inventory.CurrentCapacity,
		&inventory.CreatedAt,
		&inventory.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // 在庫が見つからない場合はnilを返す
		}
		log.Printf("プレイヤー在庫検索エラー (UserID: %s): %v", userID, err)
		return nil, err
	}

	return &inventory, nil
}

// CreateInventory は新しい在庫情報を作成します
// ユーザー登録時に呼び出されることを想定しています。
func (r *PlayerInventoryRepository) CreateInventory(inventory *models.PlayerInventory) error {
	query := `
		INSERT INTO player_inventories (id, user_id, gold, max_capacity, current_capacity, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	if inventory.ID == "" {
		inventory.ID = uuid.New().String()
	}
	now := time.Now()
	if inventory.CreatedAt.IsZero() {
		inventory.CreatedAt = now
	}
	if inventory.UpdatedAt.IsZero() {
		inventory.UpdatedAt = now
	}

	_, err := r.db.Exec(
		query,
		inventory.ID,
		inventory.UserID,
		inventory.Gold,
		inventory.MaxCapacity,
		inventory.CurrentCapacity,
		inventory.CreatedAt,
		inventory.UpdatedAt,
	)

	if err != nil {
		log.Printf("プレイヤー在庫作成エラー (UserID: %s): %v", inventory.UserID, err)
		return err
	}

	return nil
}
