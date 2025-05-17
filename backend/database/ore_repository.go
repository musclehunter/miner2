package database

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/musclehunter/miner2/models"
)

// OreRepository は鉱石関連のデータベース操作を管理します
type OreRepository struct {
	db *sql.DB
}

// NewOreRepository は新しいOreRepositoryインスタンスを作成します
func NewOreRepository(db *sql.DB) *OreRepository {
	return &OreRepository{db: db}
}

// GetAllOres は全ての鉱石情報を取得します
func (r *OreRepository) GetAllOres() ([]*models.Ore, error) {
	query := `
		SELECT id, name, rarity, purity, processing_difficulty, created_at, updated_at
		FROM ores
		ORDER BY rarity ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("鉱石一覧取得エラー: %v", err)
		return nil, err
	}
	defer rows.Close()

	var ores []*models.Ore
	for rows.Next() {
		var ore models.Ore
		err := rows.Scan(
			&ore.ID,
			&ore.Name,
			&ore.Rarity,
			&ore.Purity,
			&ore.ProcessingDifficulty,
			&ore.CreatedAt,
			&ore.UpdatedAt,
		)
		if err != nil {
			log.Printf("鉱石データスキャンエラー: %v", err)
			return nil, err
		}
		ores = append(ores, &ore)
	}

	if err = rows.Err(); err != nil {
		log.Printf("鉱石一覧読み込みエラー: %v", err)
		return nil, err
	}

	return ores, nil
}

// GetOreByID はIDから鉱石情報を取得します
func (r *OreRepository) GetOreByID(id string) (*models.Ore, error) {
	query := `
		SELECT id, name, rarity, purity, processing_difficulty, created_at, updated_at
		FROM ores
		WHERE id = ?
	`

	var ore models.Ore
	err := r.db.QueryRow(query, id).Scan(
		&ore.ID,
		&ore.Name,
		&ore.Rarity,
		&ore.Purity,
		&ore.ProcessingDifficulty,
		&ore.CreatedAt,
		&ore.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // 鉱石が見つからない場合はnilを返す
		}
		log.Printf("鉱石検索エラー: %v", err)
		return nil, err
	}

	return &ore, nil
}

// CreateOre は新しい鉱石をデータベースに保存します
func (r *OreRepository) CreateOre(ore *models.Ore) error {
	query := `
		INSERT INTO ores (id, name, rarity, purity, processing_difficulty, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	// IDが設定されていない場合は生成
	if ore.ID == "" {
		ore.ID = uuid.New().String()
	}

	now := time.Now()
	if ore.CreatedAt.IsZero() {
		ore.CreatedAt = now
	}
	if ore.UpdatedAt.IsZero() {
		ore.UpdatedAt = now
	}

	_, err := r.db.Exec(
		query,
		ore.ID,
		ore.Name,
		ore.Rarity,
		ore.Purity,
		ore.ProcessingDifficulty,
		ore.CreatedAt,
		ore.UpdatedAt,
	)
	if err != nil {
		log.Printf("鉱石作成エラー: %v", err)
		return err
	}

	return nil
}
