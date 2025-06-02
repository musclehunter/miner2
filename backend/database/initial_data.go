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

	// 初期町データ - world_setting.mdに基づく7地域
	towns := []struct {
		ID          string
		Name        string
		Description string
		PositionX   int
		PositionY   int
	}{
		{
			ID:          "1",
			Name:        "ルナフロスト",
			Description: "北方のルナリス氷原に位置する半地下都市。氷精族が住まい、常に雪に覆われた美しい都市。",
			PositionX:   50,
			PositionY:   15,
		},
		{
			ID:          "2",
			Name:        "スカイスパイア",
			Description: "北東のアエリド天空断崖に浮かぶ空中都市。羽人族が住み、浮遠石と鎖で支えられた壁々が特徴。",
			PositionX:   80,
			PositionY:   30,
		},
		{
			ID:          "3",
			Name:        "シルヴァリオン",
			Description: "東方のフェイエルフ深緑大森林に位置する螺旋都市。森エルフ族が住み、巨大な神木を中心に発展した。",
			PositionX:   85,
			PositionY:   50,
		},
		{
			ID:          "4",
			Name:        "インゴットリム",
			Description: "南東のマグノーム火山地帯に位置する鉄鋼の都市。溶炉ドワーフ族が住み、火山口に建つ鉱山都市。",
			PositionX:   70,
			PositionY:   70,
		},
		{
			ID:          "5",
			Name:        "ザル＝バディア",
			Description: "南方の広大な砂漠に位置するオアシス都市。砂漠獣人族（サンドビーストキン）が住み、地下水脈を利用した都市。",
			PositionX:   50,
			PositionY:   80,
		},
		{
			ID:          "6",
			Name:        "フォグヴェイル",
			Description: "南西の湿地に位置する霧に覆われた都市。半霊アンデッド（ミレニアル）族が住み、矢立つ石塔群で特徴づけられる。",
			PositionX:   30,
			PositionY:   70,
		},
		{
			ID:          "7",
			Name:        "キャメロス",
			Description: "西方の高原城塁地帯に位置する人間の都市。多層城壁と市場が特徴で、交易ハブとして機能している。",
			PositionX:   15,
			PositionY:   50,
		},
	}

	// 町データの挿入
	stmt, err := tx.Prepare("INSERT INTO towns (id, name, description, position_x, position_y, created_at, updated_at) VALUES (?, ?, ?, ?, ?, NOW(), NOW())")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, town := range towns {
		_, err := stmt.Exec(town.ID, town.Name, town.Description, town.PositionX, town.PositionY)
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
