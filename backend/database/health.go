package database

import (
	"fmt"
	"time"
)

// CheckDatabaseHealth はデータベース接続を確認します
func CheckDatabaseHealth() (string, error) {
	if DB == nil {
		return "", fmt.Errorf("データベース接続が初期化されていません")
	}

	// データベース接続を確認
	err := DB.Ping()
	if err != nil {
		return "", fmt.Errorf("データベース接続エラー: %v", err)
	}

	// 現在時刻をデータベースから取得してUTCで返す
	var dbTime time.Time
	query := "SELECT NOW()"
	err = DB.QueryRow(query).Scan(&dbTime)
	if err != nil {
		return "", fmt.Errorf("データベースクエリエラー: %v", err)
	}

	return fmt.Sprintf("データベース接続成功: 現在時刻 %s (UTC)", dbTime.UTC().Format(time.RFC3339)), nil
}

// CreateInitialData は初期データをデータベースに追加します
func CreateInitialData() error {
	// 初期町データを挿入
	towns := []struct {
		ID          string
		Name        string
		Description string
	}{
		{"1", "アイアンヒル", "鉄鉱石の産地として知られる古い鉱山の町。"},
		{"2", "シルバーレイク", "銀鉱石が豊富な湖のほとりにある町。"},
		{"3", "ゴールドクレスト", "金鉱脈が発見されて栄えた歴史ある町。"},
		{"4", "クリスタルヴェイル", "美しい結晶が取れる渓谷近くの町。"},
		{"5", "コッパークリーク", "銅鉱石の採掘で栄えた小さな町。"},
	}

	// トランザクション開始
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("トランザクション開始エラー: %v", err)
	}

	// 既存データの確認
	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM towns").Scan(&count)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("towns テーブル検索エラー: %v", err)
	}

	// データが既に存在する場合は何もしない
	if count > 0 {
		tx.Rollback()
		return nil
	}

	// 町データの挿入
	stmt, err := tx.Prepare("INSERT INTO towns (id, name, description, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())")
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("クエリ準備エラー: %v", err)
	}
	defer stmt.Close()

	for _, town := range towns {
		_, err := stmt.Exec(town.ID, town.Name, town.Description)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("town 挿入エラー %s: %v", town.Name, err)
		}
	}

	// コミット
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("トランザクションコミットエラー: %v", err)
	}

	return nil
}
