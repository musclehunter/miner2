package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"

	"github.com/yourusername/miner2/models"
)

// ユーザーストア（実際のプロジェクトではDBから取得）
var users = make(map[string]*models.User)

// JWTの秘密鍵（本番環境では環境変数から取得するべき）
var jwtSecret = []byte("your-secret-key")

// サインアップリクエスト
type SignupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

// ログインリクエスト
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// JWTクレーム
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// Signup_CreateMockUser はテスト用のユーザーを作成（開発用）
func Signup_CreateMockUser(email, name, password string) (*models.User, error) {
	// 新しいユーザーを作成
	user, err := models.NewUser(email, name, password)
	if err != nil {
		return nil, err
	}

	// UUIDを生成
	user.ID = uuid.New().String()

	// ユーザーを保存
	users[user.ID] = user

	return user, nil
}

// Signup はユーザー登録を処理
func Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力データが不正です", "details": err.Error()})
		return
	}

	// メールアドレスが既に使用されているかチェック
	for _, u := range users {
		if u.Email == req.Email {
			c.JSON(http.StatusConflict, gin.H{"error": "このメールアドレスは既に使用されています"})
			return
		}
	}

	// 新しいユーザーを作成
	user, err := models.NewUser(req.Email, req.Name, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー作成に失敗しました"})
		return
	}

	// UUIDを生成
	user.ID = uuid.New().String()

	// ユーザーを保存（実際のプロジェクトではDBに保存）
	users[user.ID] = user

	// JWTトークンを生成
	token, err := generateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "トークン生成に失敗しました"})
		return
	}

	// レスポンス
	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

// Login はユーザーログインを処理
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力データが不正です", "details": err.Error()})
		return
	}

	// ユーザーを検索（実際のプロジェクトではDBから検索）
	var foundUser *models.User
	for _, u := range users {
		if u.Email == req.Email {
			foundUser = u
			break
		}
	}

	// ユーザーが見つからないか、パスワードが一致しない場合
	if foundUser == nil || !foundUser.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "メールアドレスまたはパスワードが正しくありません"})
		return
	}

	// JWTトークンを生成
	token, err := generateToken(foundUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "トークン生成に失敗しました"})
		return
	}

	// レスポンス
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    foundUser.ID,
			"email": foundUser.Email,
			"name":  foundUser.Name,
		},
	})
}

// generateToken はJWTトークンを生成
func generateToken(userID string) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Printf("トークン生成エラー: %v", err)
		return "", err
	}

	return tokenString, nil
}

// AuthMiddleware は認証ミドルウェア
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "認証トークンがありません"})
			c.Abort()
			return
		}

		// "Bearer "の部分を削除
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// トークンを解析
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無効なトークンです"})
			c.Abort()
			return
		}

		// ユーザーIDをコンテキストに設定
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

// GetCurrentUser は現在のログインユーザー情報を取得
func GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証されていません"})
		return
	}

	user, exists := users[userID.(string)]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}
