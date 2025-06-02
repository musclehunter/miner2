package database

import (
	"database/sql"
	"log"
)

// MigrateTownCoordinates はtownsテーブルに座標カラムを追加し、既存データを更新します
func MigrateTownCoordinates(db *sql.DB) error {
	log.Println("町の座標情報マイグレーションを開始します...")
	
	// トランザクション開始
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	
	// カラム存在確認
	var columnExists bool
	err = tx.QueryRow(`
		SELECT COUNT(*) > 0
		FROM information_schema.columns
		WHERE table_name = 'towns'
		AND column_name = 'position_x'
	`).Scan(&columnExists)
	
	if err != nil {
		tx.Rollback()
		return err
	}
	
	// カラムが存在しない場合は追加
	if !columnExists {
		log.Println("座標カラムを追加します...")
		
		_, err = tx.Exec(`
			ALTER TABLE towns
			ADD COLUMN position_x INT NOT NULL DEFAULT 0,
			ADD COLUMN position_y INT NOT NULL DEFAULT 0
		`)
		
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		log.Println("座標カラムは既に存在します")
	}
	
	// 既存データに座標情報を更新
	townCoordinates := []struct {
		ID        string
		PositionX int
		PositionY int
	}{
		{ID: "1", PositionX: 50, PositionY: 15},  // ルナフロスト
		{ID: "2", PositionX: 80, PositionY: 30},  // スカイスパイア
		{ID: "3", PositionX: 85, PositionY: 50},  // シルヴァリオン
		{ID: "4", PositionX: 70, PositionY: 70},  // インゴットリム
		{ID: "5", PositionX: 50, PositionY: 80},  // ザル＝バディア
		{ID: "6", PositionX: 30, PositionY: 70},  // フォグヴェイル
		{ID: "7", PositionX: 15, PositionY: 50},  // キャメロス
	}
	
	stmt, err := tx.Prepare(`
		UPDATE towns
		SET position_x = ?, position_y = ?
		WHERE id = ?
	`)
	
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	
	for _, town := range townCoordinates {
		_, err := stmt.Exec(town.PositionX, town.PositionY, town.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	
	// コミット
	if err := tx.Commit(); err != nil {
		return err
	}
	
	log.Println("町の座標情報マイグレーションが完了しました")
	return nil
}
