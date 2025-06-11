package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/musclehunter/miner2/models"
)

// PlayerOreRepository はプレイヤーの鉱石所持情報に関連するデータベース操作を管理します
type PlayerOreRepository struct {
	db *sql.DB
}

// NewPlayerOreRepository は新しい PlayerOreRepository インスタンスを作成します
func NewPlayerOreRepository(db *sql.DB) *PlayerOreRepository {
	return &PlayerOreRepository{db: db}
}

// GetOresByUserID はユーザーIDで所持している鉱石一覧を取得します
func (r *PlayerOreRepository) GetOresByUserID(userID string) ([]*models.PlayerOre, error) {
	query := `
		SELECT po.id, po.user_id, po.ore_id, po.quantity, po.created_at, po.updated_at,
		       o.id, o.name, o.rarity, o.purity, o.processing_difficulty, o.created_at, o.updated_at
		FROM player_ores po
		JOIN ores o ON po.ore_id = o.id
		WHERE po.user_id = ?
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		log.Printf("プレイヤー所持鉱石一覧取得エラー (UserID: %s): %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var playerOres []*models.PlayerOre
	for rows.Next() {
		var po models.PlayerOre
		err := rows.Scan(
			&po.ID, &po.UserID, &po.OreID, &po.Quantity, &po.CreatedAt, &po.UpdatedAt,
			&po.Ore.ID, &po.Ore.Name, &po.Ore.Rarity, &po.Ore.Purity, &po.Ore.ProcessingDifficulty, &po.Ore.CreatedAt, &po.Ore.UpdatedAt,
		)
		if err != nil {
			log.Printf("プレイヤー所持鉱石データスキャンエラー: %v", err)
			return nil, err
		}
		playerOres = append(playerOres, &po)
	}

	if err = rows.Err(); err != nil {
		log.Printf("プレイヤー所持鉱石一覧読み込みエラー: %v", err)
		return nil, err
	}

	return playerOres, nil
}

// CreateOrUpdatePlayerOre はプレイヤーの鉱石所持情報を更新または作成します
func (r *PlayerOreRepository) CreateOrUpdatePlayerOre(po *models.PlayerOre) error {
	// 既存のレコードがあるか確認
	querySelect := "SELECT id, quantity FROM player_ores WHERE user_id = ? AND ore_id = ?"
	var existingID string
	var existingQuantity int
	err := r.db.QueryRow(querySelect, po.UserID, po.OreID).Scan(&existingID, &existingQuantity)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("プレイヤー所持鉱石検索エラー: %v", err)
		return err
	}

	// レコードが存在する場合 (UPDATE)
	if err != sql.ErrNoRows {
		query := "UPDATE player_ores SET quantity = ?, updated_at = ? WHERE id = ?"
		_, err = r.db.Exec(query, po.Quantity, time.Now(), existingID)
		if err != nil {
			log.Printf("プレイヤー所持鉱石更新エラー: %v", err)
			return err
		}
		return nil
	}

	// レコードが存在しない場合 (INSERT)
	queryInsert := `
		INSERT INTO player_ores (id, user_id, ore_id, quantity, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	if po.ID == "" {
		po.ID = uuid.New().String()
	}
	now := time.Now()
	if po.CreatedAt.IsZero() {
		po.CreatedAt = now
	}
	po.UpdatedAt = now

	_, err = r.db.Exec(queryInsert, po.ID, po.UserID, po.OreID, po.Quantity, po.CreatedAt, po.UpdatedAt)
	if err != nil {
		log.Printf("プレイヤー所持鉱石作成エラー: %v", err)
		return err
	}

	return nil
}
