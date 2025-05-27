package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/musclehunter/miner2/cache"
	"github.com/musclehunter/miner2/database"
	"github.com/musclehunter/miner2/models"
)

// 管理者アクセス用のシークレットキー
// 本番環境では環境変数から取得すべき
var adminSecretKey = "admin-secret-key"

// AdminAuth は管理者認証のミドルウェア
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("Authorization")
		
		// Bearer トークンの形式を想定
		if len(apiKey) > 7 && apiKey[:7] == "Bearer " {
			apiKey = apiKey[7:]
		}
		
		if apiKey != adminSecretKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "管理者権限がありません"})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// AdminLoginRequest は管理者ログインリクエスト
type AdminLoginRequest struct {
	SecretKey string `json:"secret_key" binding:"required"`
}

// AdminLogin は管理者ログイン処理を行うハンドラ
func AdminLogin(c *gin.Context) {
	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不正なリクエスト形式です"})
		return
	}
	
	// シークレットキーを検証
	if req.SecretKey != adminSecretKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "管理者キーが無効です"})
		return
	}
	
	// シンプルな実装として、シークレットキー自体をトークンとして返す
	// 本番環境では、JWTなどの安全なトークン生成方式を使用すべき
	c.JSON(http.StatusOK, gin.H{
		"token": adminSecretKey,
		"message": "管理者ログインに成功しました",
	})
}

// ユーザー管理

// GetAllUsers は全ユーザー情報を管理者向けに取得
func GetAllUsers(c *gin.Context) {
	// データベース接続確認
	if database.DB == nil {
		log.Printf("データベース接続が確立されていません")
		// データベース接続が確立されていない場合は空のリストを返す
		c.JSON(http.StatusOK, gin.H{
			"users": []map[string]interface{}{},
			"message": "データベース接続が確立されていません",
		})
		return
	}
	
	// 予期しないパニックをキャッチする
	defer func() {
		if r := recover(); r != nil {
			log.Printf("予期しないエラーが発生しました: %v", r)
			c.JSON(http.StatusOK, gin.H{
				"users": []map[string]interface{}{},
				"message": "予期しないエラーが発生しました",
			})
		}
	}()
	
	userRepo := database.NewUserRepository(database.DB)
	users, err := userRepo.GetAllUsers()
	if err != nil {
		log.Printf("ユーザー一覧取得エラー: %v", err)
		// エラーが発生した場合は空のリストを返す
		c.JSON(http.StatusOK, gin.H{
			"users": []map[string]interface{}{},
			"message": "ユーザー情報の取得に失敗しました",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

// GetUserDetail はユーザー詳細を取得
func GetUserDetail(c *gin.Context) {
	userID := c.Param("id")
	
	userRepo := database.NewUserRepository(database.DB)
	user, err := userRepo.GetUserByID(userID)
	if err != nil {
		log.Printf("ユーザー詳細取得エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー情報の取得に失敗しました"})
		return
	}
	
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// UpdateUserRequest はユーザー更新リクエスト
type UpdateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// UpdateUser はユーザー情報を更新
func UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力データが不正です"})
		return
	}
	
	// ユーザーを取得
	userRepo := database.NewUserRepository(database.DB)
	user, err := userRepo.GetUserByID(userID)
	if err != nil {
		log.Printf("ユーザー取得エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー情報の取得に失敗しました"})
		return
	}
	
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
		return
	}
	
	// 値を更新
	if req.Email != "" {
		user.Email = req.Email
	}
	
	if req.Name != "" {
		user.Name = req.Name
	}
	
	if req.Password != "" {
		if err := user.SetPassword(req.Password); err != nil {
			log.Printf("パスワード設定エラー: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "パスワードの設定に失敗しました"})
			return
		}
	}
	
	// ユーザーを更新
	userRepoUpdate := database.NewUserRepository(database.DB)
	err = userRepoUpdate.UpdateUser(user)
	if err != nil {
		log.Printf("ユーザー更新エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー情報の更新に失敗しました"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "ユーザー情報を更新しました",
		"user":    user,
	})
}

// DeleteUser はユーザーを削除
func DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	
	// ユーザーを取得して存在確認
	userRepo := database.NewUserRepository(database.DB)
	user, err := userRepo.GetUserByID(userID)
	if err != nil {
		log.Printf("ユーザー取得エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー情報の取得に失敗しました"})
		return
	}
	
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
		return
	}
	
	// ユーザーを削除
	userRepoDelete := database.NewUserRepository(database.DB)
	err = userRepoDelete.DeleteUser(userID)
	if err != nil {
		log.Printf("ユーザー削除エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザーの削除に失敗しました"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "ユーザーを削除しました",
	})
}

// 町管理

// InitAdminHandlers は管理者ハンドラを初期化
func InitAdminHandlers() {
	// 管理者用の初期化処理が必要な場合はここに追加
	log.Println("管理者ハンドラの初期化完了")
}

// GetAllTownsAdmin は全町データを管理者向けに取得
func GetAllTownsAdmin(c *gin.Context) {
	townRepo := database.NewTownRepository(database.DB)
	towns, err := townRepo.GetAllTowns()
	if err != nil {
		log.Printf("町データ取得エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "町データの取得に失敗しました"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"towns": towns,
	})
}

// UpdateTownRequest は町データ更新リクエスト
type UpdateTownRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateTown は町データを更新
func UpdateTown(c *gin.Context) {
	townID := c.Param("id")
	
	var req UpdateTownRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力データが不正です"})
		return
	}
	
	// 町データを取得
	townRepo := database.NewTownRepository(database.DB)
	town, err := townRepo.GetTownByID(townID)
	if err != nil {
		log.Printf("町データ取得エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "町データの取得に失敗しました"})
		return
	}
	
	if town == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "町が見つかりません"})
		return
	}
	
	// 値を更新
	if req.Name != "" {
		town.Name = req.Name
	}
	
	if req.Description != "" {
		town.Description = req.Description
	}
	
	// 町データを更新
	townRepoUpdate := database.NewTownRepository(database.DB)
	err = townRepoUpdate.UpdateTown(town)
	if err != nil {
		log.Printf("町データ更新エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "町データの更新に失敗しました"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "町データを更新しました",
		"town":    town,
	})
}

// CreateTownRequest は町データ作成リクエスト
type CreateTownRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// CreateTown は新しい町を作成
func CreateTown(c *gin.Context) {
	var req CreateTownRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力データが不正です"})
		return
	}
	
	// 新しい町を作成
	town := &models.Town{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
	}
	
	// 町データを保存
	townRepo := database.NewTownRepository(database.DB)
	err := townRepo.CreateTown(town)
	if err != nil {
		log.Printf("町データ作成エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "町の作成に失敗しました"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "町を作成しました",
		"town":    town,
	})
}

// DeleteTown は町を削除
func DeleteTown(c *gin.Context) {
	townID := c.Param("id")
	
	// 町データを取得して存在確認
	townRepo := database.NewTownRepository(database.DB)
	town, err := townRepo.GetTownByID(townID)
	if err != nil {
		log.Printf("町データ取得エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "町データの取得に失敗しました"})
		return
	}
	
	if town == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "町が見つかりません"})
		return
	}
	
	// 町を削除
	townRepoDelete := database.NewTownRepository(database.DB)
	err = townRepoDelete.DeleteTown(townID)
	if err != nil {
		log.Printf("町削除エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "町の削除に失敗しました"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "町を削除しました",
	})
}

// Redis管理

// GetPendingUsers は仮登録ユーザー一覧を取得
func GetPendingUsers(c *gin.Context) {
	ctx := context.Background()
	
	// 確認トークンのパターンを取得
	keys, err := cache.Client.Keys(ctx, cache.EmailVerificationPrefix+"*")
	if err != nil {
		log.Printf("仮登録ユーザー取得エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "仮登録ユーザーの取得に失敗しました"})
		return
	}
	
	var pendingUsers []gin.H
	
	// 各トークンのユーザー情報を取得
	for _, key := range keys {
		// プレフィックスを削除してトークンを取得
		token := key[len(cache.EmailVerificationPrefix):]
		
		// ユーザー情報を取得
		userJSON, err := cache.Client.Get(ctx, key)
		if err != nil {
			log.Printf("仮登録ユーザー情報取得エラー: %v", err)
			continue
		}
		
		// JSON解析
		pendingUser, err := models.PendingUserFromJSON(userJSON)
		if err != nil {
			log.Printf("JSONパースエラー: %v", err)
			continue
		}
		
		pendingUsers = append(pendingUsers, gin.H{
			"token":     token,
			"email":     pendingUser.Email,
			"name":      pendingUser.Name,
			"createdAt": pendingUser.CreatedAt,
		})
	}
	
	c.JSON(http.StatusOK, gin.H{
		"pendingUsers": pendingUsers,
	})
}

// DeletePendingUser は仮登録ユーザーを削除
func DeletePendingUser(c *gin.Context) {
	token := c.Param("token")
	ctx := context.Background()
	
	// トークンの存在確認
	exists, err := cache.Client.Exists(ctx, cache.EmailVerificationPrefix+token)
	if err != nil {
		log.Printf("トークン確認エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "トークンの確認に失敗しました"})
		return
	}
	
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "指定されたトークンが見つかりません"})
		return
	}
	
	// トークンを削除
	err = cache.DeleteEmailVerification(ctx, token)
	if err != nil {
		log.Printf("仮登録ユーザー削除エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "仮登録ユーザーの削除に失敗しました"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "仮登録ユーザーを削除しました",
	})
}

// GetRedisStatus はRedisの状態を取得
func GetRedisStatus(c *gin.Context) {
	ctx := context.Background()
	
	// キーの数を取得
	keys, err := cache.Client.Keys(ctx, "*")
	if err != nil {
		log.Printf("Redisキー取得エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Redisの状態確認に失敗しました"})
		return
	}
	
	// メール確認トークンの数
	verificationKeys, err := cache.Client.Keys(ctx, cache.EmailVerificationPrefix+"*")
	if err != nil {
		log.Printf("確認トークン取得エラー: %v", err)
		verificationKeys = []string{} // エラー時は空リストを使用
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "online",
		"stats": gin.H{
			"totalKeys":          len(keys),
			"verificationTokens": len(verificationKeys),
		},
	})
}
