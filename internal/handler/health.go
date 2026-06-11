package handler

import (
	"encoding/json"
	"net/http"
)

// Описывает формат ответа для проверки состояния сервиса
type HealthResponse struct {
	Status string `json:"status"`
}

// Обрабатывает запросы GET
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		resp := map[string]string{"error": "Method not allowed"}
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Возвращаем {"status": "ok"}
	response := HealthResponse{
		Status: "ok",
	}

	_ = json.NewEncoder(w).Encode(response)
}
