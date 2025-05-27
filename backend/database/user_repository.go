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

// DeleteUser はユーザーを削除します
func (r *UserRepository) DeleteUser(userID string) error {
	query := `
		DELETE FROM users
		WHERE id = ?
	`

	_, err := r.db.Exec(query, userID)
	if err != nil {
		log.Printf("ユーザー削除エラー: %v", err)
		return err
	}

	return nil
}

// GetAllUsers は全ユーザーを取得します
func (r *UserRepository) GetAllUsers() ([]*models.User, error) {
	query := `
		SELECT id, email, password_hash, salt, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		log.Printf("ユーザー一覧取得エラー: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		// nameフィールドはデータベースにはないがクライアント互換性のため空の値を設定
		user.Name = ""

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.Salt,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Printf("ユーザーデータスキャンエラー: %v", err)
			return nil, err
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		log.Printf("ユーザーデータ取得エラー: %v", err)
		return nil, err
	}

	return users, nil
}

// 注意: このプロジェクトの設計では、メール確認後にのみユーザーがデータベースに登録されます。
// 未確認ユーザーは全てRedisに保存されるため、このメソッドは不要です。

// 注意: このプロジェクトの設計では、確認トークンはRedisに保存され、データベースには保存されないため、
// このメソッドは不要です。
