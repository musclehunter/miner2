package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/musclehunter/miner2/cache"
	"github.com/musclehunter/miner2/database"
	"github.com/musclehunter/miner2/mail"
	"github.com/musclehunter/miner2/models"
)

// ユーザーリポジトリ
var userRepo *database.UserRepository

// メール送信者
var mailSender mail.Sender

// アプリケーションのベースURL
var baseURL string

// InitHandlersはハンドラの初期化を行います
func InitHandlers() {
	userRepo = database.NewUserRepository(database.DB)
	
	// メール送信者を初期化
	mailConfig := mail.DefaultConfig()
	
	// 開発環境ではモック送信者、本番環境ではSMTP送信者を使用
	if os.Getenv("APP_ENV") == "production" {
		mailSender = mail.NewSMTPSender(mailConfig)
		log.Println("メール送信者を初期化: SMTP送信者")
	} else {
		mailSender = mail.NewMockSender()
		log.Println("メール送信者を初期化: モック送信者")
	}
	
	// ベースURLを設定
	baseURL = os.Getenv("APP_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	
	log.Println("認証ハンドラーの初期化完了: UserRepositoryとメール送信者設定済み")
}

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
	// 既に存在するかチェック
	existingUser, err := userRepo.GetUserByEmail(email)
	if err != nil {
		log.Printf("ユーザー検索エラー: %v", err)
		return nil, err
	}

	if existingUser != nil {
		// 既に存在する場合はそのユーザーを返す
		log.Printf("ユーザーが既に存在します: %s", email)
		return existingUser, nil
	}

	// 新しいユーザーを作成
	user, err := models.NewUser(email, name, password)
	if err != nil {
		log.Printf("ユーザー作成エラー: %v", err)
		return nil, err
	}

	// UUIDを生成
	user.ID = uuid.New().String()

	// データベースに保存
	err = userRepo.CreateUser(user)
	if err != nil {
		log.Printf("ユーザー保存エラー: %v", err)
		return nil, err
	}

	log.Printf("モックユーザーを作成しました: %s", email)

	return user, nil
}

// generateVerificationToken はメール確認用のトークンを生成
func generateVerificationToken() (string, error) {
	// ランダムなトークンを生成
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	
	// Base64でエンコード
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// Signup はユーザー登録リクエストを処理
func Signup(c *gin.Context) {
	// 予期しないパニックをキャッチする
	defer func() {
		if r := recover(); r != nil {
			log.Printf("予期しないエラーが発生しました: %v", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "システムエラーが発生しました。後ほど再試行してください。"})
		}
	}()
	
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力データが不正です", "details": err.Error()})
		return
	}

	// 0. データベース接続確認
	if database.DB == nil {
		log.Printf("データベース接続が確立されていません")
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "システムが現在メンテナンス中です。後ほど再試行してください。"})
		return
	}

	// 1. データベースでメールアドレスの重複をチェック
	existingUser, err := userRepo.GetUserByEmail(req.Email)
	if err != nil {
		log.Printf("ユーザー検索エラー: %v", err)
		// 開発環境ではデータベースエラーを無視して続行
		if os.Getenv("APP_ENV") != "production" {
			log.Printf("開発環境のため、データベースエラーを無視して続行します")
			existingUser = nil
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー検索中にエラーが発生しました"})
			return
		}
	}

	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "このメールアドレスは既に使用されています"})
		return
	}

	// 2. Redisで仮登録ユーザーの重複をチェック
	ctx := context.Background()
	emailExists, err := cache.CheckEmailExists(ctx, req.Email)
	if err != nil {
		log.Printf("仮登録ユーザー検索エラー: %v", err)
		// エラーは無視して続行
	}

	if emailExists {
		c.JSON(http.StatusConflict, gin.H{"error": "このメールアドレスは確認待ちの状態です。メールを確認して登録を完了してください。"})
		return
	}
	
	// 3. 仮登録ユーザー情報を作成
	pendingUser := models.NewPendingUser(req.Email, req.Name, req.Password)
	
	// 4. 確認トークンを生成
	token, err := generateVerificationToken()
	if err != nil {
		log.Printf("トークン生成エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "確認トークンの生成に失敗しました"})
		return
	}
	
	// 5. ユーザー情報をJSON化してRedisに保存
	userJSON, err := pendingUser.ToJSON()
	if err != nil {
		log.Printf("ユーザー情報シリアライズエラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー情報の処理に失敗しました"})
		return
	}
	
	err = cache.SaveEmailVerification(ctx, token, userJSON)
	if err != nil {
		log.Printf("仮登録情報保存エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "仮登録情報の保存に失敗しました"})
		return
	}

	// 6. 確認メールを送信
	log.Printf("メール確認トークン: %s", token)
	err = mail.SendVerificationEmail(mailSender, req.Email, req.Name, token, baseURL)
	if err != nil {
		log.Printf("メール送信エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "確認メールの送信に失敗しました"})
		return
	}

	log.Printf("仮登録成功: %s", req.Email)

	// 7. 成功レスポンスを返す
	c.JSON(http.StatusOK, gin.H{
		"message": "確認メールを送信しました。メールを確認して登録を完了してください。",
		"email": req.Email,
	})
}

// Login はユーザーログインを処理
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力データが不正です", "details": err.Error()})
		return
	}

	log.Printf("ログイン試行: %s", req.Email)

	// ユーザーをデータベースから検索
	log.Printf("ユーザー検索開始: email=%s, userRepo=%v", req.Email, userRepo)
	foundUser, err := userRepo.GetUserByEmail(req.Email)
	if err != nil {
		log.Printf("ログイン時のユーザー検索エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ログイン処理中にエラーが発生しました"})
		return
	}

	// ユーザーが見つからない場合
	if foundUser == nil {
		log.Printf("ログイン失敗: ユーザーが見つかりません %s", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "メールアドレスまたはパスワードが正しくありません"})
		return
	}
	
	log.Printf("ユーザー検索結果: %+v", foundUser)

	// パスワードを検証
	if !foundUser.CheckPassword(req.Password) {
		log.Printf("ログイン失敗: パスワードが一致しません %s", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "メールアドレスまたはパスワードが正しくありません"})
		return
	}

	// JWTトークンを生成
	log.Printf("トークン生成開始: ユーザーID=%s", foundUser.ID)
	token, err := generateToken(foundUser.ID)
	if err != nil {
		log.Printf("トークン生成エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "トークン生成に失敗しました"})
		return
	}

	log.Printf("ログイン成功: メール=%s, ユーザーID=%s, トークン=%s", req.Email, foundUser.ID, token[:10]+"...")

	// レスポンス
	response := gin.H{
		"token": token,
		"user": gin.H{
			"id":    foundUser.ID,
			"email": foundUser.Email,
			"name":  foundUser.Name,
		},
	}
	log.Printf("レスポンス送信: %+v", response)
	c.JSON(http.StatusOK, response)
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

// Me は現在のログインユーザー情報を取得
// VerifyEmailRequest はメール確認リクエスト
type VerifyEmailRequest struct {
	Token string `form:"token" binding:"required"`
}

// VerifyEmail はメール確認を処理し、本登録を完了する
func VerifyEmail(c *gin.Context) {
	var req VerifyEmailRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不正なトークンです"})
		return
	}
	
	// 1. Redisからトークンに対応する仮登録情報を取得
	ctx := context.Background()
	userJSON, err := cache.GetEmailVerification(ctx, req.Token)
	if err != nil {
		log.Printf("トークン取得エラー: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "無効なトークンまたは期限切れのトークンです"})
		return
	}
	
	// 2. 仮登録情報をパース
	pendingUser, err := models.PendingUserFromJSON(userJSON)
	if err != nil {
		log.Printf("ユーザー情報パースエラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー情報の処理に失敗しました"})
		return
	}
	
	// 3. 再度メールアドレスの重複をチェック
	existingUser, err := userRepo.GetUserByEmail(pendingUser.Email)
	if err != nil {
		log.Printf("ユーザー検索エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー検索中にエラーが発生しました"})
		return
	}

	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "このメールアドレスは既に使用されています"})
		return
	}
	
	// 4. 新しいユーザーを作成
	user, err := models.NewUser(pendingUser.Email, pendingUser.Name, pendingUser.Password)
	if err != nil {
		log.Printf("ユーザー作成エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー作成に失敗しました"})
		return
	}

	// 5. UUIDを生成
	user.ID = uuid.New().String()

	// 6. ユーザーをデータベースに保存
	err = userRepo.CreateUser(user)
	if err != nil {
		log.Printf("ユーザー保存エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー保存に失敗しました"})
		return
	}
	
	// 7. Redisから仮登録情報を削除
	err = cache.DeleteEmailVerification(ctx, req.Token)
	if err != nil {
		log.Printf("仮登録情報削除エラー: %v", err)
		// 削除エラーは無視して続行
	}
	
	// 8. JWTトークンを生成
	token, err := generateToken(user.ID)
	if err != nil {
		log.Printf("トークン生成エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "トークン生成に失敗しました"})
		return
	}
	
	log.Printf("メール確認完了と本登録成功: %s", user.Email)
	
	// 9. 成功レスポンス
	c.JSON(http.StatusOK, gin.H{
		"message": "メールアドレスが確認されました。登録が完了しました。",
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

// ResendVerificationEmailRequest は確認メール再送信リクエスト
type ResendVerificationEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResendVerificationEmail は確認メールを再送信
func ResendVerificationEmail(c *gin.Context) {
	var req ResendVerificationEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力データが不正です"})
		return
	}
	
	ctx := context.Background()
	
	// 1. まずデータベースでメールアドレスが登録済みかチェック
	existingUser, err := userRepo.GetUserByEmail(req.Email)
	if err != nil {
		log.Printf("ユーザー検索エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー検索中にエラーが発生しました"})
		return
	}
	
	// 既に登録済みの場合
	if existingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "このメールアドレスは既に登録済みです。ログインしてください。"})
		return
	}
	
	// 2. Redisで仮登録情報をチェック
	emailExists, err := cache.CheckEmailExists(ctx, req.Email)
	if err != nil {
		log.Printf("仮登録ユーザー検索エラー: %v", err)
		// エラーは無視して続行
	}

	// 仮登録が見つからない場合
	if !emailExists {
		// セキュリティ上の理由からエラーを表示せず成功を裏がえる
		c.JSON(http.StatusOK, gin.H{"message": "確認メールを送信しました。メールを確認して登録を完了してください。"})
		return
	}
	
	// 3. Redis内の全トークンから対象メールに関連するトークンを探す
	// 実装の単純化のため、現在の全トークンを削除して新しいトークンを生成
	keys, err := cache.Client.Keys(ctx, cache.EmailVerificationPrefix+"*")
	if err == nil {
		for _, key := range keys {
			userJSON, err := cache.Client.Get(ctx, key)
			if err != nil {
				continue
			}
			
			pendingUser, err := models.PendingUserFromJSON(userJSON)
			if err != nil {
				continue
			}
			
			if pendingUser.Email == req.Email {
				// トークンを削除
				cache.Client.Del(ctx, key)
			}
		}
	}
	
	// 4. 新しい仮登録ユーザー情報を作成
	// ここでは完全なユーザー情報を持っていないため、簡略化された形で送信
	pendingUser := models.NewPendingUser(req.Email, "", "")

	// 5. 新しい確認トークンを生成
	token, err := generateVerificationToken()
	if err != nil {
		log.Printf("トークン生成エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "確認トークンの生成に失敗しました"})
		return
	}
	
	// 6. ユーザー情報をJSON化
	userJSON, err := pendingUser.ToJSON()
	if err != nil {
		log.Printf("ユーザー情報シリアライズエラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー情報の処理に失敗しました"})
		return
	}
	
	// 7. Redisに保存
	err = cache.SaveEmailVerification(ctx, token, userJSON)
	if err != nil {
		log.Printf("仮登録情報保存エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "仮登録情報の保存に失敗しました"})
		return
	}
	
	// 8. 確認メールを送信
	log.Printf("メール確認トークン(再送信): %s", token)
	err = mail.SendVerificationEmail(mailSender, req.Email, pendingUser.Name, token, baseURL)
	if err != nil {
		log.Printf("メール送信エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "メール送信に失敗しました"})
		return
	}
	
	// 9. 成功レスポンス
	c.JSON(http.StatusOK, gin.H{"message": "確認メールを送信しました。メールを確認して登録を完了してください。"})
}

// Me は現在のログインユーザー情報を取得
func Me(c *gin.Context) {
	// 認証ミドルウェアでセットされたユーザーIDを取得
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}
	
	// ユーザー情報を取得
	log.Printf("ユーザー情報取得: userID=%v", userID)
	user, err := userRepo.GetUserByID(userID.(string))
	if err != nil {
		log.Printf("ユーザー情報取得エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー情報の取得に失敗しました"})
		return
	}
	
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
		return
	}
	
	log.Printf("ユーザー情報取得成功: %+v", user)
	
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":           user.ID,
			"email":        user.Email,
			"name":         user.Name,
		},
	})
}
