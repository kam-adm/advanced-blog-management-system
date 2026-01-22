package handler

import (
	"advanced-blog-management-system/internal/errors/apperrors"
	"advanced-blog-management-system/internal/model"
	"advanced-blog-management-system/internal/service"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	userService service.UserServiceInterface
}

// NewAuthHandler создает новый экземпляр AuthHandler
// TODO: Инициализировать с userService (интерфейс)
func NewAuthHandler(userService service.UserServiceInterface) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

// Register обрабатывает POST /api/register
// TODO: Проверить метод, распарсить JSON, валидировать, вызвать userService.Register()
// Вернуть TokenResponse со статусом 201 или ошибку с нужным кодом
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать
}

// Login обрабатывает POST /api/login
// TODO: Проверить метод, распарсить JSON, валидировать, вызвать userService.Login()
// Вернуть TokenResponse со статусом 200 или ошибку с нужным кодом
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать
}

// GetProfile получает профиль текущего пользователя
// TODO: Проверить метод GET, получить userID из контекста, вернуть UserResponse (опционально)
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать получение профиля (опционально)
}

// respondWithJSON отправляет JSON ответ с заданным статус кодом
// TODO: Установить Content-Type, WriteHeader, закодировать JSON
func (h *AuthHandler) respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	// TODO: реализовать
}

// respondWithError отправляет JSON ошибку
// TODO: Создать ErrorResponse и отправить используя respondWithJSON()
func (h *AuthHandler) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	// TODO: реализовать
}
