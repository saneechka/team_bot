package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"team_bot/config"
	"team_bot/internal/bot"
	"team_bot/internal/repository/sqlrepo"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(cfg.Database.Path), 0755); err != nil {
		log.Fatalf("Error creating data directory: %v", err)
	}

	db, err := sql.Open("sqlite3", cfg.Database.Path)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	authRepo := sqlrepo.NewAuthRepository(db)

	botAPI, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}
	botAPI.Debug = cfg.Bot.Debug

	telegramBot := bot.New(botAPI, authRepo, cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Received shutdown signal")
		cancel()
	}()

	telegramBot.Start(ctx)
}
