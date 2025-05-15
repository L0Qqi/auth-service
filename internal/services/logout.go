package services

import "database/sql"

func InvalidateAllUserTokens(db *sql.DB, userID string) error {
	_, err := db.Exec(`UPDATE refresh_tokens SET online = false WHERE user_id = $1`, userID)
	return err
}
