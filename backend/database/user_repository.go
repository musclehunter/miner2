package database

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/musclehunter/miner2/models"
)

// UserRepository はユーザー関連のデータベース操作を管理します
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository は新しいUserRepositoryインスタンスを作成します
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser は新しいユーザーをデータベースに保存します
func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, salt, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	// IDが設定されていない場合は生成
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = now
	}

	_, err := r.db.Exec(
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.Salt,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		log.Printf("ユーザー作成エラー: %v", err)
		return err
	}

	return nil
}

// GetUserByEmail はメールアドレスからユーザーを検索します
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, salt, created_at, updated_at
		FROM users
		WHERE email = ?
	`

	var user models.User
	// nameフィールドはデータベースにはないがクライアント互換性のため空の値を設定
	user.Name = ""

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Salt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // ユーザーが見つからない場合はnilを返す
		}
		log.Printf("ユーザー検索エラー: %v", err)
		return nil, err
	}

	return &user, nil
}

// GetUserByID はIDからユーザーを検索します
func (r *UserRepository) GetUserByID(id string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, salt, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	var user models.User
	// nameフィールドはデータベースにはないがクライアント互換性のため空の値を設定
	user.Name = ""

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Salt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // ユーザーが見つからない場合はnilを返す
		}
		log.Printf("ユーザー検索エラー: %v", err)
		return nil, err
	}

	return &user, nil
}

// UpdateUser はユーザー情報を更新します
func (r *UserRepository) UpdateUser(user *models.User) error {
	query := `
		UPDATE users
		SET email = ?, password_hash = ?, salt = ?, updated_at = ?
		WHERE id = ?
	`

	user.UpdatedAt = time.Now()

	_, err := r.db.Exec(
		query,
		user.Email,
		user.PasswordHash,
		user.Salt,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		log.Printf("ユーザー更新エラー: %v", err)
		return err
	}

	return nil
}
