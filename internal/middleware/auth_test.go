package middleware

import (
	"advanced-blog-management-system/pkg/auth"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Проверяет защиту эндпоинтов от несанкционированного доступа
func TestAuthMiddleware_RequireAuth(t *testing.T) {
	secret := "middleware-test-secret-key"
	manager := auth.NewJWTManager(secret, 1)
	mw := NewAuthMiddleware(manager)

	// Создаем тестовый токен
	validToken, _, _ := manager.GenerateToken(1, "user@mail.ru", "user")

	// Простая заглушка-хендлер, которая возвращает 200 OK, если управление дошло до неё
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "Valid bearer token provided",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Missing auth header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Malformed token header format",
			authHeader:     "MalformedTokenString123",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Invalid or expired token content",
			authHeader:     "Bearer wrong-token-content",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/posts", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rr := httptest.NewRecorder()

			// Запускаем RequireAuth, обернутый в адаптер
			handlerToTest := mw.RequireAuth(nextHandler)
			handlerToTest.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, tt.expectedStatus)
			}
		})
	}
}
