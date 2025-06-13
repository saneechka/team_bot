package handler

import (
	"context"
	"fmt"
	"log"
	"time"

	"team_bot/internal/model"
	"team_bot/internal/repository/sqlrepo"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AuthHandler struct {
	bot        *tgbotapi.BotAPI
	repo       *sqlrepo.AuthRepository
	adminUsers []string
}

func NewAuthHandler(bot *tgbotapi.BotAPI, repo *sqlrepo.AuthRepository, adminUsers []string) *AuthHandler {
	return &AuthHandler{
		bot:        bot,
		repo:       repo,
		adminUsers: adminUsers,
	}
}

func (h *AuthHandler) CheckAdminAccess(ctx context.Context, userID int64, chatID int64) (bool, error) {
	isAdmin, err := h.repo.IsAdmin(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("error checking admin status: %v", err)
	}

	if !isAdmin {
		msg := tgbotapi.NewMessage(chatID, "❌ Доступ запрещён. У вас нет прав администратора.")
		if _, err := h.bot.Send(msg); err != nil {
			log.Printf("Error sending access denied message: %v", err)
		}
		return false, nil
	}

	return true, nil
}

func (h *AuthHandler) HandleStart(ctx context.Context, update *tgbotapi.Update) error {
	// Проверяем, является ли пользователь администратором по username
	isAdmin := false
	username := update.Message.From.UserName
	for _, adminUsername := range h.adminUsers {
		if username == adminUsername {
			isAdmin = true
			break
		}
	}

	user := &model.User{
		ID:          update.Message.From.ID,
		Username:    username,
		ChatID:      update.Message.Chat.ID,
		CreatedTime: time.Now(),
		IsAdmin:     isAdmin,
	}

	if err := h.repo.SaveUser(ctx, user); err != nil {
		log.Printf("Error saving user: %v", err)
		return fmt.Errorf("error saving user: %v", err)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Say Hello", "hello_btn"),
		),
	)

	// Добавляем информацию о статусе администратора в приветственное сообщение
	adminStatus := ""
	if isAdmin {
		adminStatus = "\n✅ Вы зарегистрированы как администратор."
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID,
		fmt.Sprintf("Привет, %s! Я бот для управления командой.%s", username, adminStatus))
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ReplyMarkup = keyboard

	if _, err := h.bot.Send(msg); err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}

	return nil
}

func (h *AuthHandler) HandleAdmin(ctx context.Context, update *tgbotapi.Update) error {
	isAdmin, err := h.repo.IsAdmin(ctx, update.Message.From.ID)
	if err != nil {
		return fmt.Errorf("error checking admin status: %v", err)
	}

	var responseText string
	if isAdmin {
		responseText = "✅ Вы являетесь администратором."
	} else {
		responseText = "❌ Вы не являетесь администратором."
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseText)
	if _, err := h.bot.Send(msg); err != nil {
		return fmt.Errorf("error sending admin status message: %v", err)
	}

	return nil
}
