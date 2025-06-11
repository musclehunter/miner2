package database

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/musclehunter/miner2/models"
)

// CreatePlayerBase は新しいプレイヤー拠点をデータベースに作成します。
// この関数は、必要であればトランザクション内で player_inventories も作成します。
func CreatePlayerBase(tx *sql.Tx, userID, townID string) (*models.PlayerBase, error) {
	// ユーザーのインベントリが既に存在するか確認
	var existingInventoryID string
	err := tx.QueryRow("SELECT id FROM player_inventories WHERE user_id = ?", userID).Scan(&existingInventoryID)

	// インベントリが存在しない場合のみ作成
	if err == sql.ErrNoRows {
		// プレイヤーインベントリの初期化
		newInventory := &models.PlayerInventory{
			ID:              uuid.New().String(),
			UserID:          userID,
			Gold:            100, // 初期所持金
			MaxCapacity:     100, // 初期最大容量
			CurrentCapacity: 0,
		}

		queryInventory := `
			INSERT INTO player_inventories (id, user_id, gold, max_capacity, current_capacity)
			VALUES (?, ?, ?, ?, ?)
		`
		_, err = tx.Exec(queryInventory, newInventory.ID, newInventory.UserID, newInventory.Gold, newInventory.MaxCapacity, newInventory.CurrentCapacity)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		// その他のデータベースエラー
		return nil, err
	}

	// 新しい拠点の作成
	newBase := &models.PlayerBase{
		ID:     uuid.New().String(),
		UserID: userID,
		TownID: townID,
		Level:  1,
	}

	queryBase := `
		INSERT INTO player_bases (id, user_id, town_id, level)
		VALUES (?, ?, ?, ?)
	`
	_, err = tx.Exec(queryBase, newBase.ID, newBase.UserID, newBase.TownID, newBase.Level)
	if err != nil {
		return nil, err
	}

	return newBase, nil
}

// GetAllPlayerBases はすべてのプレイヤー拠点をデータベースから取得します。
func GetAllPlayerBases(db *sql.DB) ([]*models.PlayerBase, error) {
	query := `
		SELECT id, user_id, town_id, level, created_at, updated_at
		FROM player_bases
		ORDER BY created_at DESC
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bases []*models.PlayerBase
	for rows.Next() {
		base := &models.PlayerBase{}
		err := rows.Scan(&base.ID, &base.UserID, &base.TownID, &base.Level, &base.CreatedAt, &base.UpdatedAt)
		if err != nil {
			return nil, err
		}
		bases = append(bases, base)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bases, nil
}
