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
	// ログ出力
	log.Printf("CreateBase called with request body: %v", c.Request.Body)
	var req CreateBaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// ログ出力
	log.Printf("createbase request json: %v", c)

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// ログ出力
	log.Printf("CreateBase userID: %v", userID)

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// ログ出力
	log.Printf("CreateBase database begin")

	newBase, err := database.CreateBaseWithWarehouse(tx, userID.(string), req.TownID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create base"})
		log.Printf("CreateBase database create base with warehouse error: %v", err)
		return
	}

	// ログ出力
	log.Printf("CreateBase create player base")

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		log.Printf("CreateBase database commit error: %v", err)
		return
	}

	// ログ出力
	log.Printf("CreateBase database commit")

	c.JSON(http.StatusCreated, newBase)
}

// GetAllBasesHandler retrieves all bases
func GetAllBasesHandler(c *gin.Context) {
	bases, err := database.GetAllBases(db)
	if err != nil {
		log.Printf("データベースからの拠点取得に失敗しました: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bases"})
		return
	}

	// 拠点が一件もない場合はnullではなく空の配列を返す
	if len(bases) == 0 {
		c.JSON(http.StatusOK, []*models.Base{})
		return
	}

	c.JSON(http.StatusOK, bases)
}
