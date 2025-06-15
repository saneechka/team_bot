package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"team_bot/internal/model"
	"team_bot/internal/repository/sqlrepo"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AuthService struct {
	bot        *tgbotapi.BotAPI
	repo       *sqlrepo.AuthRepository
	adminUsers []string
}

func NewAuthService(bot *tgbotapi.BotAPI, repo *sqlrepo.AuthRepository, adminUsers []string) *AuthService {
	return &AuthService{
		bot:        bot,
		repo:       repo,
		adminUsers: adminUsers,
	}
}

func (s *AuthService) CheckAdminAccess(ctx context.Context, userID int64, chatID int64) (bool, error) {
	isAdmin, err := s.repo.IsAdmin(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("error checking admin status: %v", err)
	}

	if !isAdmin {
		msg := tgbotapi.NewMessage(chatID, "❌ Доступ запрещён. У вас нет прав администратора.")
		if _, err := s.bot.Send(msg); err != nil {
			log.Printf("Error sending access denied message: %v", err)
		}
		return false, nil
	}

	return true, nil
}

func (s *AuthService) IsUserAdmin(username string) bool {
	for _, adminUsername := range s.adminUsers {
		if username == adminUsername {
			return true
		}
	}
	return false
}

func (s *AuthService) CreateUser(ctx context.Context, userID int64, username string, chatID int64, isAdmin bool) (*model.User, error) {
	user := &model.User{
		ID:          userID,
		Username:    username,
		Name:   "", 
		Surname:    "", 
		ChatID:      chatID,
		CreatedTime: time.Now(),
		IsAdmin:     isAdmin,
	}

	if err := s.repo.SaveUser(ctx, user); err != nil {
		return nil, fmt.Errorf("error saving user: %v", err)
	}

	return user, nil
}

func (s *AuthService) GetUserByID(ctx context.Context, userID int64) (*model.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}


func (s *AuthService) UpdatePersonalInfo(ctx context.Context, userID int64, firstName, lastName string) error {
	return s.repo.AddPersonalInfo(ctx, userID, firstName, lastName)
}


func (s *AuthService) UpdateUserPersonalInfo(ctx context.Context, user *model.User) error {
	return s.repo.UpdatePersonalInfo(ctx, user)
}
