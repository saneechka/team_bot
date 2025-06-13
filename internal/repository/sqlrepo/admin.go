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



func (r *AuthRepository) CreateInviteToken(ctx context.Context, token *model.InviteToken) error {
	query := `
		INSERT INTO invite_tokens (token, created_by, created_at, expires_at, is_active, usage_count, max_usage)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.ExecContext(ctx, query,
		token.Token,
		token.CreatedBy,
		token.CreatedAt,
		token.ExpiresAt,
		token.IsActive,
		token.UsageCount,
		token.MaxUsage,
	)
	if err != nil {
		return fmt.Errorf("error creating invite token: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %v", err)
	}
	token.ID = id

	return nil
}

func (r *AuthRepository) GetActiveInviteToken(ctx context.Context) (*model.InviteToken, error) {
	query := `
		SELECT id, token, created_by, created_at, expires_at, is_active, usage_count, max_usage
		FROM invite_tokens
		WHERE is_active = true AND expires_at > datetime('now')
		ORDER BY created_at DESC
		LIMIT 1
	`
	var token model.InviteToken
	err := r.db.QueryRowContext(ctx, query).Scan(
		&token.ID,
		&token.Token,
		&token.CreatedBy,
		&token.CreatedAt,
		&token.ExpiresAt,
		&token.IsActive,
		&token.UsageCount,
		&token.MaxUsage,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting active invite token: %v", err)
	}
	return &token, nil
}

func (r *AuthRepository) GetInviteTokenByToken(ctx context.Context, tokenStr string) (*model.InviteToken, error) {
	query := `
		SELECT id, token, created_by, created_at, expires_at, is_active, usage_count, max_usage
		FROM invite_tokens
		WHERE token = ?
	`
	var token model.InviteToken
	err := r.db.QueryRowContext(ctx, query, tokenStr).Scan(
		&token.ID,
		&token.Token,
		&token.CreatedBy,
		&token.CreatedAt,
		&token.ExpiresAt,
		&token.IsActive,
		&token.UsageCount,
		&token.MaxUsage,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting invite token: %v", err)
	}
	return &token, nil
}

func (r *AuthRepository) DeactivateAllInviteTokens(ctx context.Context) error {
	query := `UPDATE invite_tokens SET is_active = false`
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("error deactivating invite tokens: %v", err)
	}
	return nil
}

func (r *AuthRepository) IncrementTokenUsage(ctx context.Context, tokenID int64) error {
	query := `UPDATE invite_tokens SET usage_count = usage_count + 1 WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, tokenID)
	if err != nil {
		return fmt.Errorf("error incrementing token usage: %v", err)
	}
	return nil
}

func (r *AuthRepository) UserExists(ctx context.Context, userID int64) (bool, error) {
	query := `SELECT 1 FROM users WHERE id = ? LIMIT 1`
	var exists int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("error checking user existence: %v", err)
	}
	return true, nil
}
