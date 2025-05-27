package database

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/musclehunter/miner2/models"
)

// TownRepository は町関連のデータベース操作を管理します
type TownRepository struct {
	db *sql.DB
}

// NewTownRepository は新しいTownRepositoryインスタンスを作成します
func NewTownRepository(db *sql.DB) *TownRepository {
	return &TownRepository{db: db}
}

// GetAllTowns は全ての町情報を取得します
func (r *TownRepository) GetAllTowns() ([]*models.Town, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM towns
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("町一覧取得エラー: %v", err)
		return nil, err
	}
	defer rows.Close()

	var towns []*models.Town
	for rows.Next() {
		var town models.Town
		err := rows.Scan(
			&town.ID,
			&town.Name,
			&town.Description,
			&town.CreatedAt,
			&town.UpdatedAt,
		)
		if err != nil {
			log.Printf("町データスキャンエラー: %v", err)
			return nil, err
		}
		towns = append(towns, &town)
	}

	if err = rows.Err(); err != nil {
		log.Printf("町一覧読み込みエラー: %v", err)
		return nil, err
	}

	return towns, nil
}

// GetTownByID はIDから町情報を取得します
func (r *TownRepository) GetTownByID(id string) (*models.Town, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM towns
		WHERE id = ?
	`

	var town models.Town
	err := r.db.QueryRow(query, id).Scan(
		&town.ID,
		&town.Name,
		&town.Description,
		&town.CreatedAt,
		&town.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // 町が見つからない場合はnilを返す
		}
		log.Printf("町検索エラー: %v", err)
		return nil, err
	}

	return &town, nil
}

// CreateTown は新しい町をデータベースに保存します
func (r *TownRepository) CreateTown(town *models.Town) error {
	query := `
		INSERT INTO towns (id, name, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	// IDが設定されていない場合は生成
	if town.ID == "" {
		town.ID = uuid.New().String()
	}

	now := time.Now()
	if town.CreatedAt.IsZero() {
		town.CreatedAt = now
	}
	if town.UpdatedAt.IsZero() {
		town.UpdatedAt = now
	}

	_, err := r.db.Exec(
		query,
		town.ID,
		town.Name,
		town.Description,
		town.CreatedAt,
		town.UpdatedAt,
	)
	if err != nil {
		log.Printf("町作成エラー: %v", err)
		return err
	}

	return nil
}

// UpdateTown は町の情報を更新します
func (r *TownRepository) UpdateTown(town *models.Town) error {
	query := `
		UPDATE towns
		SET name = ?, description = ?, updated_at = ?
		WHERE id = ?
	`

	town.UpdatedAt = time.Now()

	_, err := r.db.Exec(
		query,
		town.Name,
		town.Description,
		town.UpdatedAt,
		town.ID,
	)
	if err != nil {
		log.Printf("町更新エラー: %v", err)
		return err
	}

	return nil
}

// DeleteTown は町を削除します
func (r *TownRepository) DeleteTown(townID string) error {
	query := `
		DELETE FROM towns
		WHERE id = ?
	`

	_, err := r.db.Exec(query, townID)
	if err != nil {
		log.Printf("町削除エラー: %v", err)
		return err
	}

	return nil
}
