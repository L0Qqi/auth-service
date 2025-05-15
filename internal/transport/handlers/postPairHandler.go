package handlers

import (
	"auth-service/internal/app/jwt"
	"auth-service/internal/app/refresh"
	"auth-service/internal/app/tokens"
	"database/sql"
	"encoding/json"
	"net/http"
)

// HandleGetTokens выдает новую пару access/refresh токенов по user_id
// @Summary Получение токенов
// @Description Возвращает новую пару access и refresh токенов по user_id
// @Tags auth
// @Accept json
// @Produce json
// @Param user_id query string true "Идентификатор пользователя (GUID)"
// @Success 200 {object} tokens.TokenPair
// @Failure 400 {string} string "missing user_id"
// @Failure 500 {string} string "failed to generate token"
// @Router /tokens [get]
func HandleGetTokens(db *sql.DB, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("user_id")
		if userID == "" {
			http.Error(w, "missing user_id", http.StatusBadRequest)
			return
		}

		accessToken, err := jwt.GenerateAccessToken(userID, []byte(jwtSecret))
		if err != nil {
			http.Error(w, "failed to generate access token", http.StatusInternalServerError)
			return
		}

		refreshToken, err := refresh.IssueRefreshToken(db, userID, r)
		if err != nil {
			http.Error(w, "failed to issue refresh token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokens.TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	}
}
