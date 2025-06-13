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


type InviteToken struct {
	ID         int64     `json:"id"`
	Token      string    `json:"token"`
	CreatedBy  int64     `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	IsActive   bool      `json:"is_active"`
	UsageCount int       `json:"usage_count"`
	MaxUsage   int       `json:"max_usage"`
}


func IsTelegramAdmin(username string, adminUsernames []string) bool {
	for _, adminUsername := range adminUsernames {
		if username == adminUsername {
			return true
		}
	}
	return false
}
