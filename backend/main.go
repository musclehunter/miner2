package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	
	"github.com/yourusername/miner2/handlers"
)

// モックユーザーの作成（開発用）
func createMockUser() {
	// テスト用ユーザーを作成
	mockUser, _ := handlers.Signup_CreateMockUser("test@example.com", "テストユーザー", "password123")
	log.Printf("モックユーザーを作成しました: %s", mockUser.Email)
}

func main() {
	// 開発用のモックユーザーを作成
	createMockUser()
	// Ginモードの設定（開発環境はデバッグモード）
	gin.SetMode(gin.DebugMode)

	// ルーターの初期化
	r := gin.Default()

	// CORSの設定
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 基本的なルート
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "採掘人 API サーバーが稼働中です",
		})
	})

	// APIグループ
	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "健康状態は良好です！",
			})
		})

		// 認証関連のエンドポイント
		auth := api.Group("/auth")
		{
			auth.POST("/signup", handlers.Signup)
			auth.POST("/login", handlers.Login)
			auth.GET("/me", handlers.AuthMiddleware(), handlers.GetCurrentUser)
		}

		// 認証が必要なエンドポイント
		secured := api.Group("/")
		secured.Use(handlers.AuthMiddleware())
		{
			// 今後ここに認証が必要なエンドポイントを追加
			// 例: プレイヤーデータ、ゲーム状態、拠点情報など
		}
	}

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // デフォルトポート
	}

	serverAddr := fmt.Sprintf(":%s", port)
	log.Printf("サーバーを起動しています: %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}
