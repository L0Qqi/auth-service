package refresh

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Генерация токена
func generateSecureToken() (string, error) {
	s := make([]byte, 32)
	_, err := rand.Read(s)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(s), nil
}

// Хеширование токена
func hashToken(token string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func getIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// Генерация и сохранение refresh токена
func IssueRefreshToken(db *sql.DB, userID string, r *http.Request) (string, error) {
	token, err := generateSecureToken()
	if err != nil {
		return "", err
	}

	tokenHash, err := hashToken(token)
	if err != nil {
		return "", err
	}

	id := uuid.New().String()
	createdAt := time.Now()
	finishAt := createdAt.Add(7 * 24 * time.Hour) // 7 дней
	ip := getIP(r)

	_, err = db.Exec(`
        INSERT INTO refresh_tokens (id, user_id, token_hash, revoked, created_at, finish_at, ip)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`, id, userID, tokenHash, false, createdAt, finishAt, ip)

	if err != nil {
		return "", fmt.Errorf("не удалось сохранить токен в базе данных: %w", err)
	}

	return token, nil
}
