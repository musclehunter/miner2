package database

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/musclehunter/miner2/models"
)

// CreateBaseWithWarehouse は新しい拠点と倉庫をトランザクション内で作成します。
func CreateBaseWithWarehouse(tx *sql.Tx, userID, townID string) (*models.Base, error) {
	// 新しい拠点の作成
	newBase := &models.Base{
		ID:     uuid.New().String(),
		UserID: userID,
		TownID: townID,
	}

	queryBase := `
		INSERT INTO bases (id, user_id, town_id)
		VALUES (?, ?, ?)
	`
	_, err := tx.Exec(queryBase, newBase.ID, newBase.UserID, newBase.TownID)
	if err != nil {
		return nil, err
	}

	// 新しい倉庫の作成
	newWarehouse := &models.Warehouse{
		ID:       uuid.New().String(),
		BaseID:   newBase.ID,
		Level:    1,
		Capacity: 100, // 初期容量
	}

	queryWarehouse := `
		INSERT INTO warehouses (id, base_id, level, capacity)
		VALUES (?, ?, ?, ?)
	`
	_, err = tx.Exec(queryWarehouse, newWarehouse.ID, newWarehouse.BaseID, newWarehouse.Level, newWarehouse.Capacity)
	if err != nil {
		return nil, err
	}

	return newBase, nil
}

// GetAllBases はすべての拠点をデータベースから取得します。
func GetAllBases(db *sql.DB) ([]*models.Base, error) {
	query := `
		SELECT id, user_id, town_id, created_at, updated_at
		FROM bases
		ORDER BY created_at DESC
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bases []*models.Base
	for rows.Next() {
		base := &models.Base{}
		err := rows.Scan(&base.ID, &base.UserID, &base.TownID, &base.CreatedAt, &base.UpdatedAt)
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
