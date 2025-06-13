package model

import "slices"

import "time"


type Message struct {
	ID        int64  `json:"id"`
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	Timestamp time.Time  `json:"timestamp"`
	Type      string `json:"type"`
}


type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	ChatID    int64     `json:"chat_id"`
	IsAdmin   bool      `json:"is_admin"`
}


func IsTelegramAdmin(username string, adminUsernames []string) bool {
	return slices.Contains(adminUsernames, username)
}

