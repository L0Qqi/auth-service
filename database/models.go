package db

import "time"

// RefreshToken структура для хранения refresh токена
type RefreshToken struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	TokenHash string    `json:"token_hash"`
	Revoked   bool      `json:"revoked"`
	CreatedAt time.Time `json:"created_at"`
	FinishAt  time.Time `json:"finish_at"`
}
