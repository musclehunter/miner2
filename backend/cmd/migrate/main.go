package main

import (
	"log"
	"os"

	"github.com/musclehunter/miner2/database"
)

func main() {
	log.Println("データベースマイグレーションツールを開始します...")
	
	// データベース接続
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("データベース接続エラー: %v", err)
		os.Exit(1)
	}
	defer db.Close()
	
	// 町の座標情報マイグレーション実行
	if err := database.MigrateTownCoordinates(db); err != nil {
		log.Fatalf("マイグレーション実行エラー: %v", err)
		os.Exit(1)
	}
	
	log.Println("マイグレーションが正常に完了しました")
}
