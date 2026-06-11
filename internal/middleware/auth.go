package middleware

import (
	"advanced-blog-management-system/pkg/auth"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

// contextKey - кастомный тип, предотвращающий коллизии ключей в контексте запроса
type contextKey string

const (
	UserIDKey    contextKey = "userID"
	UserEmailKey contextKey = "userEmail"
	UserNameKey  contextKey = "username"
)

// Управляет авторизацией с использованием JWTManager
type AuthMiddleware struct {
	jwtManager *auth.JWTManager
}

// Создает новый экземпляр AuthMiddleware
func NewAuthMiddleware(jwtManager *auth.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

// Жестко блокирует запросы без валидного JWT токена (для POST/PUT/DELETE методов)
func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := extractToken(r)
		if tokenStr == "" {
			respondWithError(w, "Missing or malformed authorization token", http.StatusUnauthorized)
			return
		}

		// Валидируем токен через метод объекта JWTManager
		claims, err := m.jwtManager.ValidateToken(tokenStr)
		if err != nil {
			respondWithError(w, "Invalid or expired token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Наполняем контекст всеми данными из токена
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
		ctx = context.WithValue(ctx, UserNameKey, claims.Username)

		next(w, r.WithContext(ctx))
	}
}

// Извлекает данные пользователя, если токен передан, но не блокирует гостей блога
func (m *AuthMiddleware) OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := extractToken(r)
		if tokenStr == "" {
			next(w, r)
			return
		}

		claims, err := m.jwtManager.ValidateToken(tokenStr)
		if err != nil {
			// Токен испорчен или просрочен — игнорируем его и пропускаем пользователя как гостя
			next(w, r)
			return
		}

		// Если токен валидный — добавляем информацию в контекст
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
		ctx = context.WithValue(ctx, UserNameKey, claims.Username)

		next(w, r.WithContext(ctx))
	}
}

// Безопасно извлекает userID из контекста
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDKey).(int)
	return userID, ok
}

// Безопасно извлекает email из контекста
func GetUserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(UserEmailKey).(string)
	return email, ok
}

// Безопасно извлекает username из контекста
func GetUsernameFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(UserNameKey).(string)
	return username, ok
}

// Парсит заголовок и возвращает чистую строку JWT токена
func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}

	return parts[1]
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, message string, code int) {
	type ErrorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, ErrorResponse{Error: message})
}
