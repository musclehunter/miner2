package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/musclehunter/miner2/database"
)

// ゲーム関連のリポジトリ
var (
	townRepo *database.TownRepository
	oreRepo  *database.OreRepository
)

// InitGameHandlers はゲームハンドラーの初期化を行います
func InitGameHandlers() {
	townRepo = database.NewTownRepository(database.DB)
	oreRepo = database.NewOreRepository(database.DB)
}

// GetAllTowns は全ての町情報を返すハンドラーです
func GetAllTowns(c *gin.Context) {
	towns, err := townRepo.GetAllTowns()
	if err != nil {
		log.Printf("町一覧取得エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "町の情報取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"towns": towns,
	})
}

// GetTownByID は指定IDの町情報を返すハンドラーです
func GetTownByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "町IDが指定されていません"})
		return
	}

	town, err := townRepo.GetTownByID(id)
	if err != nil {
		log.Printf("町検索エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "町の情報取得に失敗しました"})
		return
	}

	if town == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "指定された町が見つかりません"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"town": town,
	})
}

// GetAllOres は全ての鉱石情報を返すハンドラーです
func GetAllOres(c *gin.Context) {
	ores, err := oreRepo.GetAllOres()
	if err != nil {
		log.Printf("鉱石一覧取得エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "鉱石の情報取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ores": ores,
	})
}

// GetOreByID は指定IDの鉱石情報を返すハンドラーです
func GetOreByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "鉱石IDが指定されていません"})
		return
	}

	ore, err := oreRepo.GetOreByID(id)
	if err != nil {
		log.Printf("鉱石検索エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "鉱石の情報取得に失敗しました"})
		return
	}

	if ore == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "指定された鉱石が見つかりません"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ore": ore,
	})
}
