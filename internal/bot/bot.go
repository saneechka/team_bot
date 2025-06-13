package bot

import (
	"context"
	_ "fmt"
	"log"

	"team_bot/config"
	"team_bot/internal/handler"
	"team_bot/internal/repository/sqlrepo"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api         *tgbotapi.BotAPI
	authRepo    *sqlrepo.AuthRepository
	authHandler *handler.AuthHandler
	config      *config.Config
}

func New(api *tgbotapi.BotAPI, authRepo *sqlrepo.AuthRepository, cfg *config.Config) *Bot {
	return &Bot{
		api:         api,
		authRepo:    authRepo,
		authHandler: handler.NewAuthHandler(api, authRepo, cfg.TelegramAdmins.Usernames),
		config:      cfg,
	}
}

func (b *Bot) Start(ctx context.Context) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := b.api.GetUpdatesChan(updateConfig)
	log.Println("Bot started successfully")

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down bot...")
			return
		case update := <-updates:
			if update.Message != nil {
				b.handleMessage(ctx, &update)
			}
		}
	}
}

func (b *Bot) handleMessage(ctx context.Context, update *tgbotapi.Update) {
	message := update.Message

	if message.IsCommand() && message.Command() == "start" {
		if err := b.authHandler.HandleStart(ctx, update); err != nil {
			log.Printf("Error handling /start: %v", err)
		}
		return
	}

	if message.IsCommand() && message.Command() == "admin" {
		if err := b.authHandler.HandleAdmin(ctx, update); err != nil {
			log.Printf("Error handling /admin: %v", err)
		}
		return
	}
	if !b.checkAdminAccess(ctx, message.From.ID, message.Chat.ID) {
		return
	}
}

func (b *Bot) checkAdminAccess(ctx context.Context, userID int64, chatID int64) bool {
	isAdmin, err := b.authRepo.IsAdmin(ctx, userID)
	if err != nil {
		log.Printf("Error checking admin status: %v", err)
		return false
	}

	if !isAdmin {
		msg := tgbotapi.NewMessage(chatID, "❌ Доступ запрещён. Эта команда доступна только администраторам.")
		if _, err := b.api.Send(msg); err != nil {
			log.Printf("Error sending access denied message: %v", err)
		}
		return false
	}

	return true
}
