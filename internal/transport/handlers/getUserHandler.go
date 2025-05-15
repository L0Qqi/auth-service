package handlers

import (
	"auth-service/internal/app/jwt"
	"encoding/json"
	"net/http"
	"strings"
)

// HandleGetUser возвращает GUID текущего пользователя
// @Summary Получение userID
// @Tags auth
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 401 {string} string "unauthorized"
// @Router /me [get]
// @Security BearerAuth
func HandleGetUser(jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "missing or invalid authorization header", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := jwt.ParseAccessToken(token, []byte(jwtSecret))
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"user_id": userID})
	}
}
