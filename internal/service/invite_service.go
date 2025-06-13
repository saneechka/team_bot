package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"team_bot/internal/model"
	"team_bot/internal/repository/sqlrepo"
)

type InviteService struct {
	repo *sqlrepo.AuthRepository
}

func NewInviteService(repo *sqlrepo.AuthRepository) *InviteService {
	return &InviteService{
		repo: repo,
	}
}

func (s *InviteService) GenerateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("error generating random bytes: %v", err)
	}
	return hex.EncodeToString(bytes), nil
}

func (s *InviteService) CreateInviteLink(ctx context.Context, createdBy int64, durationHours int, maxUsage int) (*model.InviteToken, error) {

	if err := s.repo.DeactivateAllInviteTokens(ctx); err != nil {
		return nil, fmt.Errorf("error deactivating existing tokens: %v", err)
	}

	tokenStr, err := s.GenerateToken()
	if err != nil {
		return nil, fmt.Errorf("error generating token: %v", err)
	}

	token := &model.InviteToken{
		Token:      tokenStr,
		CreatedBy:  createdBy,
		CreatedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(time.Duration(durationHours) * time.Hour),
		IsActive:   true,
		UsageCount: 0,
		MaxUsage:   maxUsage,
	}

	if err := s.repo.CreateInviteToken(ctx, token); err != nil {
		return nil, fmt.Errorf("error creating invite token: %v", err)
	}

	return token, nil
}

func (s *InviteService) GetInviteLink(ctx context.Context) (*model.InviteToken, error) {
	return s.repo.GetActiveInviteToken(ctx)
}

func (s *InviteService) ValidateAndUseToken(ctx context.Context, tokenStr string) (*model.InviteToken, error) {

	token, err := s.repo.GetInviteTokenByToken(ctx, tokenStr)
	if err != nil {
		return nil, fmt.Errorf("error getting token: %v", err)
	}

	if token == nil {
		return nil, fmt.Errorf("токен не найден")
	}

	if !token.IsActive {
		return nil, fmt.Errorf("токен неактивен")
	}

	if time.Now().After(token.ExpiresAt) {
		return nil, fmt.Errorf("срок действия токена истек")
	}

	if token.UsageCount >= token.MaxUsage {
		return nil, fmt.Errorf("достигнут лимит использований токена")
	}

	if err := s.repo.IncrementTokenUsage(ctx, token.ID); err != nil {
		return nil, fmt.Errorf("error incrementing token usage: %v", err)
	}

	token.UsageCount++

	return token, nil
}

func (s *InviteService) FormatInviteLink(botUsername, token string) string {
	return fmt.Sprintf("https://t.me/%s?start=%s", botUsername, token)
}
