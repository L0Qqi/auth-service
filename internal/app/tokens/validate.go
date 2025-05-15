package tokens

import (
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func ValidateRefreshToken(db *sql.DB, userID string, token string, userAgent string, ip string) error {
	var storedHash string
	var tokenID int
	var revoked bool

	query := `SELECT id, token_hash, revoked FROM refresh_tokens WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1`
	err := db.QueryRow(query, userID).Scan(&tokenID, &storedHash, &revoked)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("токен не найден для пользователя: %w", err)
		}
		return err
	}
	if revoked {
		return fmt.Errorf("токен уже был использован: %w", err)
	}

	// Сравниваем токен и хеш
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(token))
	if err != nil {
		return fmt.Errorf("токен недействителен: %w", err)
	}

	if _, err := db.Exec(`UPDATE refresh_tokens SET revoked = true WHERE id = $1 AND token_hash = $2`, tokenID, storedHash); err != nil {
		return fmt.Errorf("не удалось обновить статус токена: %w", err)
	}

	return nil
}
