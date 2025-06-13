package handler

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"team_bot/internal/repository/sqlrepo"
)


type AuthHandler struct {
	bot  *tgbotapi.BotAPI
	repo *sqlrepo.AuthRepository
}


func NewAuthHandler(bot *tgbotapi.BotAPI, repo *sqlrepo.AuthRepository) *AuthHandler {
	return &AuthHandler{
		bot:  bot,
		repo: repo,
	}
}


func (h *AuthHandler) HandleStart(ctx context.Context, update *tgbotapi.Update) error {
	user := &sqlrepo.User{
		ID:        update.Message.From.ID,
		Username:  update.Message.From.UserName,
		ChatID:    update.Message.Chat.ID,
		
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

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, 
		fmt.Sprintf("Привет, %s! Я бот для управления командой.", update.Message.From.UserName))
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ReplyMarkup = keyboard

	if _, err := h.bot.Send(msg); err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}

	return nil
}


func (h *AuthHandler) HandleHelloButton(ctx context.Context, callback *tgbotapi.CallbackQuery) error {
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "hello")
	msg.ReplyToMessageID = callback.Message.MessageID
	if _, err := h.bot.Send(msg); err != nil {
		return fmt.Errorf("error sending hello message: %v", err)
	}

	callbackResp := tgbotapi.NewCallback(callback.ID, "")
	_, _ = h.bot.Request(callbackResp)
	return nil
}


func (h *AuthHandler) HandleAdmin(ctx context.Context, update *tgbotapi.Update) error {
	isAdmin, err := h.repo.IsAdmin(ctx, update.Message.From.ID)
	if err != nil {
		return fmt.Errorf("error checking admin status: %v", err)
	}

	if !isAdmin {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "У вас нет прав администратора.")
		if _, err := h.bot.Send(msg); err != nil {
			return fmt.Errorf("error sending message: %v", err)
		}
		return nil
	}


	switch update.Message.Command() {
	case "addadmin":
		return h.handleAddAdmin(ctx, update)
	case "removeadmin":
		return h.handleRemoveAdmin(ctx, update)
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда администратора.")
		if _, err := h.bot.Send(msg); err != nil {
			return fmt.Errorf("error sending message: %v", err)
		}
	}

	return nil
}


func (h *AuthHandler) handleAddAdmin(ctx context.Context, update *tgbotapi.Update) error {

	args := update.Message.CommandArguments()
	if args == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Укажите ID пользователя.")
		if _, err := h.bot.Send(msg); err != nil {
			return fmt.Errorf("error sending message: %v", err)
		}
		return nil
	}


	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Функция в разработке.")
	if _, err := h.bot.Send(msg); err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}

	return nil
}


func (h *AuthHandler) handleRemoveAdmin(ctx context.Context, update *tgbotapi.Update) error {

	args := update.Message.CommandArguments()
	if args == "" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Укажите ID пользователя.")
		if _, err := h.bot.Send(msg); err != nil {
			return fmt.Errorf("error sending message: %v", err)
		}
		return nil
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Функция в разработке.")
	if _, err := h.bot.Send(msg); err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}

	return nil
} 