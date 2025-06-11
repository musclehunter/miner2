package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/musclehunter/miner2/database"
)

var (
	playerInventoryRepo *database.PlayerInventoryRepository
	playerOreRepo       *database.PlayerOreRepository
	playerItemRepo      *database.PlayerItemRepository
)

func InitInventoryHandlers(db *sql.DB) {
	playerInventoryRepo = database.NewPlayerInventoryRepository(db)
	playerOreRepo = database.NewPlayerOreRepository(db)
	playerItemRepo = database.NewPlayerItemRepository(db)
}

// GetMyInventory はログイン中のユーザーの在庫情報を取得します
func GetMyInventory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ユーザーが認証されていません"})
		return
	}

	inventory, err := playerInventoryRepo.GetInventoryByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "在庫情報の取得に失敗しました"})
		return
	}

	ores, err := playerOreRepo.GetOresByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "所持鉱石の取得に失敗しました"})
		return
	}

	items, err := playerItemRepo.GetItemsByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "所持アイテムの取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"inventory": inventory,
		"ores":      ores,
		"items":     items,
	})
}
