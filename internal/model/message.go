package model

import (
	"time"
)

type Message struct {
	ID        int64     `json:"id"`
	ChatID    int64     `json:"chat_id"`
	Text      string    `json:"text"`
	UserID    int64     `json:"user_id"`
	Username  string    `json:"username"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
}

type User struct {
	ID          int64     `json:"id"`
	Username    string    `json:"username"`
	ChatID      int64     `json:"chat_id"`
	CreatedTime time.Time `json:"created_time"`
	IsAdmin     bool      `json:"is_admin"`
}

// IsTelegramAdmin проверяет, является ли пользователь администратором
// на основе списка административных username
func IsTelegramAdmin(username string, adminUsernames []string) bool {
	for _, adminUsername := range adminUsernames {
		if username == adminUsername {
			return true
		}
	}
	return false
}
