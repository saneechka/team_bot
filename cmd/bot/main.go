package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	"team_bot/config"
	"team_bot/internal/handler"
	"team_bot/internal/repository/sqlrepo"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Create data directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(cfg.Database.Path), 0755); err != nil {
		log.Fatalf("Error creating data directory: %v", err)
	}

	// Initialize SQLite database connection
	db, err := sql.Open("sqlite3", cfg.Database.Path)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Initialize AuthRepository
	authRepo := sqlrepo.NewAuthRepository(db)

	// Initialize bot
	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	bot.Debug = cfg.Bot.Debug

	// Set up AuthHandler
	authHandler := handler.NewAuthHandler(bot, authRepo)

	// Set up update configuration
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	// Create context that will be canceled on SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Received shutdown signal")
		cancel()
	}()

	// Start receiving updates
	updates := bot.GetUpdatesChan(updateConfig)

	log.Println("Bot started successfully")

	// Process updates
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down bot...")
			return
		case update := <-updates:
			if update.Message != nil {
				// Handle /start command
				if update.Message.IsCommand() && update.Message.Command() == "start" {
					if err := authHandler.HandleStart(ctx, &update); err != nil {
						log.Printf("Error handling /start: %v", err)
					}
					continue
				}

				// Handle /admin command
				if update.Message.IsCommand() && update.Message.Command() == "admin" {
					if err := authHandler.HandleAdmin(ctx, &update); err != nil {
						log.Printf("Error handling /admin: %v", err)
					}
					continue
				}
			}

			// Обработка callback-запросов от кнопок
			if update.CallbackQuery != nil {
				if update.CallbackQuery.Data == "hello_btn" {
					if err := authHandler.HandleHelloButton(ctx, update.CallbackQuery); err != nil {
						log.Printf("Error handling hello button: %v", err)
					}
				}
			}
		}
	}
} 