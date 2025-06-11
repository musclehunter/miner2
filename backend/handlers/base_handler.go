package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/musclehunter/miner2/database"
	"github.com/musclehunter/miner2/models"
)

// CreateBaseRequest は拠点作成リクエストのボディを表します。
type CreateBaseRequest struct {
	TownID string `json:"town_id" binding:"required"`
}

// CreateBase は新しい拠点を作成するハンドラーです。
func CreateBase(c *gin.Context) {
	var req CreateBaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	newBase, err := database.CreatePlayerBase(tx, userID.(string), req.TownID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create base"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, newBase)
}

// GetAllBasesHandler retrieves all player bases
func GetAllBasesHandler(c *gin.Context) {
	bases, err := database.GetAllPlayerBases(database.DB)
	if err != nil {
		log.Printf("データベースからの拠点取得に失敗しました: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bases"})
		return
	}

	// 拠点が一件もない場合はnullではなく空の配列を返す
	if len(bases) == 0 {
		c.JSON(http.StatusOK, []*models.PlayerBase{})
		return
	}

	c.JSON(http.StatusOK, bases)
}
