package sqlrepo

import (
	"context"
	"database/sql"
	"fmt"

	"team_bot/internal/model"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) SaveUser(ctx context.Context, user *model.User) error {
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
		user.CreatedTime,
		user.IsAdmin,
	)
	if err != nil {
		return fmt.Errorf("error saving user: %v", err)
	}
	return nil
}

func (r *AuthRepository) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT id, username, chat_id, created_at, is_admin
		FROM users
		WHERE id = ?
	`
	var user model.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.ChatID,
		&user.CreatedTime,
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

func (r *AuthRepository) GetUserByChatID(ctx context.Context, chatID int64) (*model.User, error) {
	query := `
		SELECT id, username, chat_id, created_at, is_admin
		FROM users
		WHERE chat_id = ?
	`
	var user model.User
	err := r.db.QueryRowContext(ctx, query, chatID).Scan(
		&user.ID,
		&user.Username,
		&user.ChatID,
		&user.CreatedTime,
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


func (r *AuthRepository) SetAdminStatus(ctx context.Context, userID int64, isAdmin bool) error {
	query := `UPDATE users SET is_admin = ? WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, isAdmin, userID)
	if err != nil {
		return fmt.Errorf("error setting admin status: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", userID)
	}

	return nil
}


func (r *AuthRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
		SELECT id, username, chat_id, created_at, is_admin
		FROM users
		WHERE username = ?
	`
	var user model.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.ChatID,
		&user.CreatedTime,
		&user.IsAdmin,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user by username: %v", err)
	}
	return &user, nil
}
