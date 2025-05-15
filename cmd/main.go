// @title Auth Service API
// @version 1.0
// @description Документация API для сервиса аутентификации
// @host localhost:8080
// @BasePath /

package main

import (
	db "auth-service/database"
	_ "auth-service/docs"
	"auth-service/internal/transport/handlers"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer dbConn.Close()

	jwtSecret := os.Getenv("JWT_SECRET")
	webhookURL := os.Getenv("WEBHOOK_URL")

	http.Handle("/swagger/", httpSwagger.WrapHandler)

	http.HandleFunc("/api/tokens", handlers.HandleGetTokens(dbConn, jwtSecret))
	http.HandleFunc("/api/tokens/refresh", handlers.HandleRefreshTokens(dbConn, jwtSecret, webhookURL))
	http.HandleFunc("/api/user", handlers.HandleGetUser(jwtSecret))
	http.HandleFunc("/api/logout", handlers.HandleLogout(dbConn, jwtSecret))

	port := os.Getenv("AUTH_PORT")
	if port == "" {
		port = strconv.Itoa(8080)
	}
	addr := ":" + port

	log.Printf("Сервер запущен на %s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}

}
