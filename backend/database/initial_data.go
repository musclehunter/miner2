package database

import (
	"log"
)

// CreateInitialTowns は初期の町データを作成します
func CreateInitialTowns() error {
	// データベースハンドルのチェック
	if DB == nil {
		return ErrDatabaseNotInitialized
	}

	// トランザクション開始
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	// 既存データの確認
	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM towns").Scan(&count)
	if err != nil {
		tx.Rollback()
		return err
	}

	// データが既に存在する場合は何もしない
	if count > 0 {
		tx.Rollback()
		return nil
	}

	// 初期町データ
	towns := []struct {
		ID          string
		Name        string
		Description string
	}{
		{
			ID:          "1",
			Name:        "アイアンヒル",
			Description: "鉄鉱石の産地として知られる古い鉱山の町。丘陵地帯に位置し、町の周りには多くの鉄鉱山が点在しています。",
		},
		{
			ID:          "2",
			Name:        "シルバーレイク",
			Description: "銀鉱石が豊富な湖のほとりにある町。美しい湖の底には伝説の銀脈があると言われています。",
		},
		{
			ID:          "3",
			Name:        "ゴールドクレスト",
			Description: "金鉱脈が発見されて栄えた歴史ある町。多くの採掘人が富を求めてこの地に集まります。",
		},
		{
			ID:          "4",
			Name:        "クリスタルヴェイル",
			Description: "美しい結晶が取れる渓谷近くの町。透明度の高い水晶や希少な宝石類が採掘できることで有名です。",
		},
		{
			ID:          "5",
			Name:        "コッパークリーク",
			Description: "銅鉱石の採掘で栄えた小さな町。小川のほとりに位置し、周辺の山々には多くの銅鉱脈が走っています。",
		},
	}

	// 町データの挿入
	stmt, err := tx.Prepare("INSERT INTO towns (id, name, description, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, town := range towns {
		_, err := stmt.Exec(town.ID, town.Name, town.Description)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// コミット
	if err := tx.Commit(); err != nil {
		return err
	}

	log.Println("初期町データを作成しました")
	return nil
}

// CreateInitialOres は初期の鉱石データを作成します
func CreateInitialOres() error {
	// データベースハンドルのチェック
	if DB == nil {
		return ErrDatabaseNotInitialized
	}

	// トランザクション開始
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	// 既存データの確認
	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM ores").Scan(&count)
	if err != nil {
		tx.Rollback()
		return err
	}

	// データが既に存在する場合は何もしない
	if count > 0 {
		tx.Rollback()
		return nil
	}

	// 初期鉱石データ
	ores := []struct {
		ID                  string
		Name                string
		Rarity              int
		Purity              float64
		ProcessingDifficulty int
	}{
		{
			ID:                  "1",
			Name:                "石炭",
			Rarity:              1,
			Purity:              1.0,
			ProcessingDifficulty: 1,
		},
		{
			ID:                  "2",
			Name:                "鉄鉱石",
			Rarity:              1,
			Purity:              1.0,
			ProcessingDifficulty: 2,
		},
		{
			ID:                  "3",
			Name:                "銅鉱石",
			Rarity:              2,
			Purity:              1.0,
			ProcessingDifficulty: 3,
		},
		{
			ID:                  "4",
			Name:                "銀鉱石",
			Rarity:              3,
			Purity:              0.9,
			ProcessingDifficulty: 4,
		},
		{
			ID:                  "5",
			Name:                "金鉱石",
			Rarity:              4,
			Purity:              0.8,
			ProcessingDifficulty: 5,
		},
		{
			ID:                  "6",
			Name:                "ダイヤモンド原石",
			Rarity:              5,
			Purity:              0.7,
			ProcessingDifficulty: 7,
		},
		{
			ID:                  "7",
			Name:                "エメラルド原石",
			Rarity:              5,
			Purity:              0.7,
			ProcessingDifficulty: 6,
		},
		{
			ID:                  "8",
			Name:                "サファイア原石",
			Rarity:              5,
			Purity:              0.7,
			ProcessingDifficulty: 6,
		},
		{
			ID:                  "9",
			Name:                "ルビー原石",
			Rarity:              5,
			Purity:              0.7,
			ProcessingDifficulty: 6,
		},
		{
			ID:                  "10",
			Name:                "ミスリル鉱石",
			Rarity:              6,
			Purity:              0.5,
			ProcessingDifficulty: 10,
		},
	}

	// 鉱石データの挿入
	stmt, err := tx.Prepare("INSERT INTO ores (id, name, rarity, purity, processing_difficulty, created_at, updated_at) VALUES (?, ?, ?, ?, ?, NOW(), NOW())")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, ore := range ores {
		_, err := stmt.Exec(ore.ID, ore.Name, ore.Rarity, ore.Purity, ore.ProcessingDifficulty)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// コミット
	if err := tx.Commit(); err != nil {
		return err
	}

	log.Println("初期鉱石データを作成しました")
	return nil
}
