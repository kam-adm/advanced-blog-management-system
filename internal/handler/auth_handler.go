package handler

import (
	"advanced-blog-management-system/internal/errors/apperrors"
	"advanced-blog-management-system/internal/middleware"
	"advanced-blog-management-system/internal/model"
	"advanced-blog-management-system/internal/service"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	userService service.UserServiceInterface
}

// Создает новый экземпляр AuthHandler
func NewAuthHandler(userService service.UserServiceInterface) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

// Обрабатывает POST /api/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input model.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Вызов встроенной валидации структуры
	if err := input.Validate(); err != nil {
		HandleServiceError(w, err)
		return
	}

	// Сервис теперь возвращает указатель на структуру model.TokenResponse
	tokenResponse, err := h.userService.Register(r.Context(), &input)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	h.respondWithJSON(w, tokenResponse, http.StatusCreated)
}

// Обрабатывает POST /api/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input model.UserLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Вызов встроенной валидации структуры
	if err := input.Validate(); err != nil {
		HandleServiceError(w, err)
		return
	}

	// Сервис теперь возвращает указатель на структуру model.TokenResponse
	tokenResponse, err := h.userService.Login(r.Context(), &input)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	h.respondWithJSON(w, tokenResponse, http.StatusOK)
}

// Получает профиль текущего пользователя
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Извлекаем userID из middleware
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok || userID == 0 {
		WriteError(w, apperrors.ErrUnauthorized.Error(), http.StatusUnauthorized)
		return
	}

	h.respondWithJSON(w, map[string]int{"user_id": userID}, http.StatusOK)
}

// Отправляет JSON ответ с заданным статус кодом
func (h *AuthHandler) respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
