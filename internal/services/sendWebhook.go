package services

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func SendWebhook(url, userID, ip, userAgent string) {
	body := map[string]string{
		"user_id":    userID,
		"ip":         ip,
		"user_agent": userAgent,
	}
	jsonData, _ := json.Marshal(body)
	_, _ = http.Post(url, "application/json", bytes.NewBuffer(jsonData))
}
