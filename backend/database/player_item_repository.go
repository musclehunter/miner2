package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/musclehunter/miner2/models"
)

// PlayerItemRepository はプレイヤーのアイテム所持情報に関連するデータベース操作を管理します
type PlayerItemRepository struct {
	db *sql.DB
}

// NewPlayerItemRepository は新しい PlayerItemRepository インスタンスを作成します
func NewPlayerItemRepository(db *sql.DB) *PlayerItemRepository {
	return &PlayerItemRepository{db: db}
}

// GetItemsByUserID はユーザーIDで所持しているアイテム一覧を取得します
func (r *PlayerItemRepository) GetItemsByUserID(userID string) ([]*models.PlayerItem, error) {
	query := `
		SELECT pi.id, pi.user_id, pi.item_id, pi.quantity, pi.created_at, pi.updated_at,
		       i.id, i.name, i.rarity, i.description, i.created_at, i.updated_at
		FROM player_items pi
		JOIN items i ON pi.item_id = i.id
		WHERE pi.user_id = ?
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		log.Printf("プレイヤー所持アイテム一覧取得エラー (UserID: %s): %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var playerItems []*models.PlayerItem
	for rows.Next() {
		var pi models.PlayerItem
		err := rows.Scan(
			&pi.ID, &pi.UserID, &pi.ItemID, &pi.Quantity, &pi.CreatedAt, &pi.UpdatedAt,
			&pi.Item.ID, &pi.Item.Name, &pi.Item.Rarity, &pi.Item.Description, &pi.Item.CreatedAt, &pi.Item.UpdatedAt,
		)
		if err != nil {
			log.Printf("プレイヤー所持アイテムデータスキャンエラー: %v", err)
			return nil, err
		}
		playerItems = append(playerItems, &pi)
	}

	if err = rows.Err(); err != nil {
		log.Printf("プレイヤー所持アイテム一覧読み込みエラー: %v", err)
		return nil, err
	}

	return playerItems, nil
}

// CreateOrUpdatePlayerItem はプレイヤーのアイテム所持情報を更新または作成します
func (r *PlayerItemRepository) CreateOrUpdatePlayerItem(pi *models.PlayerItem) error {
	// 既存のレコードがあるか確認
	querySelect := "SELECT id, quantity FROM player_items WHERE user_id = ? AND item_id = ?"
	var existingID string
	var existingQuantity int
	err := r.db.QueryRow(querySelect, pi.UserID, pi.ItemID).Scan(&existingID, &existingQuantity)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("プレイヤー所持アイテム検索エラー: %v", err)
		return err
	}

	// レコードが存在する場合 (UPDATE)
	if err != sql.ErrNoRows {
		query := "UPDATE player_items SET quantity = ?, updated_at = ? WHERE id = ?"
		_, err = r.db.Exec(query, pi.Quantity, time.Now(), existingID)
		if err != nil {
			log.Printf("プレイヤー所持アイテム更新エラー: %v", err)
			return err
		}
		return nil
	}

	// レコードが存在しない場合 (INSERT)
	queryInsert := `
		INSERT INTO player_items (id, user_id, item_id, quantity, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	if pi.ID == "" {
		pi.ID = uuid.New().String()
	}
	now := time.Now()
	if pi.CreatedAt.IsZero() {
		pi.CreatedAt = now
	}
	pi.UpdatedAt = now

	_, err = r.db.Exec(queryInsert, pi.ID, pi.UserID, pi.ItemID, pi.Quantity, pi.CreatedAt, pi.UpdatedAt)
	if err != nil {
		log.Printf("プレイヤー所持アイテム作成エラー: %v", err)
		return err
	}

	return nil
}
