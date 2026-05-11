package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

type healthResponse struct {
	Status  int    `json:"status"`
	Service string `json:"service"`
	Time    string `json:"time"`
}

// Health godoc
// @Summary Health check
// @Description Проверка сервиса
// @Tags system
// @Produce json
// @Success 200 {object} healthResponse
// @Router /health [get]
func Health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(healthResponse{
		Status:  http.StatusOK,
		Service: "auth",
		Time:    time.Now().UTC().Format(time.RFC3339),
	})
}
