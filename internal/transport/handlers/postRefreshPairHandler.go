package handlers

import (
	"auth-service/internal/app/jwt"
	"auth-service/internal/app/refresh"
	"auth-service/internal/app/tokens"
	"auth-service/internal/services"
	"database/sql"
	"encoding/json"
	"net"
	"net/http"
)

type RefreshRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// HandleRefreshTokens обновляет пару токенов
// @Summary Обновление токенов
// @Description Обновляет access и refresh токены при валидной паре
// @Tags auth
// @Accept json
// @Produce json
// @Param body body RefreshRequest true "Refresh tokens"
// @Success 200 {object} tokens.TokenPair
// @Failure 400 {string} string "invalid request"
// @Failure 401 {string} string "unauthorized"
// @Router /tokens/refresh [post]
func HandleRefreshTokens(db *sql.DB, jwtSecret string, webhookURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RefreshRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		userID, err := jwt.ParseAccessToken(req.AccessToken, []byte(jwtSecret))
		if err != nil {
			http.Error(w, "invalid access token", http.StatusUnauthorized)
			return
		}

		userAgent := r.UserAgent()
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		err = tokens.ValidateRefreshToken(db, userID, req.RefreshToken, userAgent, ip)
		if err != nil {
			if err.Error() == "user-agent has changed" {
				_ = services.InvalidateAllUserTokens(db, userID)
				http.Error(w, "user-agent changed, deauthorized", http.StatusUnauthorized)
				return
			}
			if err.Error() == "refresh token has already been used" {
				_ = services.InvalidateAllUserTokens(db, userID)
				http.Error(w, "refresh token reuse detected", http.StatusUnauthorized)
				return
			}
			if err.Error() == "new IP address detected" {
				go services.SendWebhook(webhookURL, userID, ip, userAgent)
			}
		}

		newAccess, err := jwt.GenerateAccessToken(userID, []byte(jwtSecret))
		if err != nil {
			http.Error(w, "failed to generate access token", http.StatusInternalServerError)
			return
		}

		newRefresh, err := refresh.IssueRefreshToken(db, userID, r)
		if err != nil {
			http.Error(w, "failed to issue refresh token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokens.TokenPair{
			AccessToken:  newAccess,
			RefreshToken: newRefresh,
		})
	}
}
