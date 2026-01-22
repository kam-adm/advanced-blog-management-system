package middleware

import (
	"advanced-blog-management-system/pkg/auth"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

// contextKey - кастомный тип, чтобы избежать коллизии
type contextKey string

const (
	// TODO: Определить константы для ключей контекста
	// UserIDKey, UserEmailKey, UserNameKey типа contextKey
	UserIDKey    contextKey = "userID"
	UserEmailKey contextKey = "userEmail"
	UserNameKey  contextKey = "username"
)

// TODO: Добавить структура AuthMiddleware с полем jwtManager
// Конструктор NewAuthMiddleware(jwtManager *auth.JWTManager) *AuthMiddleware
type AuthMiddleware struct {
	jwtManager *auth.JWTManager
}

// TODO: Реализовать конструктор NewAuthMiddleware
func NewAuthMiddleware(jwtManager *auth.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

// TODO: Реализовать метод RequireAuth(next http.HandlerFunc) http.HandlerFunc
// Извлечь токен из заголовка Authorization, валидировать,
// добавить userID, email, username в контекст, передать управление дальше
// Вернуть 401 если токена нет или он невалидный
func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: реализовать
		next(w, r)
	}
}

// TODO: Реализовать метод OptionalAuth(next http.HandlerFunc) http.HandlerFunc
// Попытаться извлечь и валидировать токен, но не требовать его
// Если токена нет или он невалидный - пропустить и передать управление дальше
// Если токен валидный - добавить userID, email, username в контекст
func (m *AuthMiddleware) OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: реализовать
		next(w, r)
	}
}

// TODO: Реализовать GetUserIDFromContext(ctx context.Context)
// Получить значение из контекста по UserIDKey и привести к типу int
// Вернуть (userID int, ok bool)
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	// TODO: реализовать
	return 0, false
}

// TODO: Реализовать GetUserEmailFromContext(ctx context.Context)
// Получить значение из контекста по UserEmailKey и привести к типу string
func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	// TODO: реализовать
	return "", false
}

// TODO: Реализовать GetUsernameFromContext(ctx context.Context)
// Получить значение из контекста по UserNameKey и привести к типу string
func GetUsernameFromContext(ctx context.Context) (string, bool) {
	// TODO: реализовать
	return "", false
}

// TODO: Реализовать приватную функцию extractToken(r *http.Request)
// Получить заголовок Authorization, проверить формат "Bearer <token>"
// Вернуть токен (вторую часть) или пустую строку если формат неверный
func extractToken(r *http.Request) string {
	// TODO: реализовать
	return ""
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, message string, code int) {
	type ErrorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, ErrorResponse{Error: message})
}
