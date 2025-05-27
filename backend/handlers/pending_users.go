package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/musclehunter/miner2/cache"
)

// GetAllPendingUsers は未確認ユーザー（仮登録ユーザー）の一覧を取得するハンドラー
func GetAllPendingUsers(c *gin.Context) {
	ctx := context.Background()
	
	// Redisからすべての未確認ユーザーを取得
	pendingUsers, err := cache.GetAllPendingUsers(ctx)
	if err != nil {
		log.Printf("未確認ユーザー一覧取得エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "未確認ユーザー情報の取得に失敗しました",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"pendingUsers": pendingUsers,
	})
}
