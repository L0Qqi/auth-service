package handlers

import (
	"auth-service/internal/app/jwt"
	"auth-service/internal/services"
	"database/sql"
	"net/http"
	"strings"
)

// HandleLogout деавторизует пользователя
// @Summary Выход пользователя
// @Tags auth
// @Produce plain
// @Success 200 {string} string "logged out successfully"
// @Failure 401 {string} string "invalid token"
// @Failure 500 {string} string "failed to logout"
// @Router /logout [post]
// @Security BearerAuth
func HandleLogout(db *sql.DB, jwtSecret string) http.HandlerFunc {
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

		err = services.InvalidateAllUserTokens(db, userID)
		if err != nil {
			http.Error(w, "failed to logout", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("logged out successfully"))
	}
}
