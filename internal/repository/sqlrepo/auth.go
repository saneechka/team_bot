package sqlrepo

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type AuthRepository struct {
	db *sql.DB
}


func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}


type User struct {
	ID        int64
	Username  string
	ChatID    int64
	CreatedAt time.Time
	IsAdmin   bool
}


func (r *AuthRepository) SaveUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, username, chat_id, created_at, is_admin)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			username = excluded.username,
			chat_id = excluded.chat_id,
			is_admin = excluded.is_admin
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Username,
		user.ChatID,
		user.CreatedAt,
		user.IsAdmin,
	)
	if err != nil {
		return fmt.Errorf("error saving user: %v", err)
	}
	return nil
}


func (r *AuthRepository) GetUserByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, username, chat_id, created_at, is_admin
		FROM users
		WHERE id = ?
	`
	var user User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.ChatID,
		&user.CreatedAt,
		&user.IsAdmin,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}
	return &user, nil
}


func (r *AuthRepository) GetUserByChatID(ctx context.Context, chatID int64) (*User, error) {
	query := `
		SELECT id, username, chat_id, created_at, is_admin
		FROM users
		WHERE chat_id = ?
	`
	var user User
	err := r.db.QueryRowContext(ctx, query, chatID).Scan(
		&user.ID,
		&user.Username,
		&user.ChatID,
		&user.CreatedAt,
		&user.IsAdmin,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}
	return &user, nil
}

func (r *AuthRepository) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	query := `SELECT is_admin FROM users WHERE id = ?`
	var isAdmin bool
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&isAdmin)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("error checking admin status: %v", err)
	}
	return isAdmin, nil
} 