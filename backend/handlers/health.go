package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/musclehunter/miner2/database"
)

// HealthCheck はシステムの健全性を確認するハンドラー
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "健康状態は良好です！",
	})
}

// DatabaseHealthCheck はデータベース接続の健全性を確認するハンドラー
func DatabaseHealthCheck(c *gin.Context) {
	message, err := database.CheckDatabaseHealth()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "エラー",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": message,
	})
}
