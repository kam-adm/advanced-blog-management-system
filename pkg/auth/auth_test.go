package auth

import (
	"testing"
	"time"
)

// Проверяет полный цикл выпуска и валидации JWT токена
func TestJWTManager_Lifecycle(t *testing.T) {
	secret := "super-secret-key-for-testing-purposes"
	manager := NewJWTManager(secret, 1) // Токен на 1 час

	userID := 42
	email := "test@mail.ru"
	username := "user"

	// 1. Тест генерации
	tokenStr, expiresAt, err := manager.GenerateToken(userID, email, username)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	if tokenStr == "" {
		t.Error("generated token string is empty")
	}
	if expiresAt.Before(time.Now()) {
		t.Error("token expiration time is in the past")
	}

	// 2. Тест валидации легитимного токена
	claims, err := manager.ValidateToken(tokenStr)
	if err != nil {
		t.Fatalf("failed to validate authentic token: %v", err)
	}
	if claims.UserID != userID || claims.Email != email || claims.Username != username {
		t.Errorf("claims data mismatch: got %d/%s/%s", claims.UserID, claims.Email, claims.Username)
	}

	// 3. Тест валидации скомпрометированного токена
	brokenManager := NewJWTManager("wrong-secret-key", 1)
	_, err = brokenManager.ValidateToken(tokenStr)
	if err == nil {
		t.Error("expected error when validating token with incorrect secret, got nil")
	}
}

// Проверяет криптографическое хеширование и верификацию паролей
func TestPassword_HashingAndVerification(t *testing.T) {
	password := "Q1we3r4r5t6!"

	// 1. Тест успешного хеширования
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	if hash == password {
		t.Error("password was not hashed, hash matches plain text")
	}

	// 2. Тест пустых данных
	_, err = HashPassword("")
	if err == nil {
		t.Error("expected error for empty password, got nil")
	}

	// 3. Тест успешного сравнения
	if !CheckPassword(password, hash) {
		t.Error("password verification failed for authentic plain text")
	}

	// 4. Тест неверного пароля
	if CheckPassword("WrongPassword123", hash) {
		t.Error("verification passed for incorrect password text")
	}
}

// Проверяет работу правил сложности паролей
func TestPassword_StrengthValidation(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"Valid strong password", "Q1we3r4r5t6!", false},
		{"Too short", "Q1w!", true},
		{"No uppercase", "q1we3r4r5t6!", true},
		{"No lowercase", "Q1WE3R4R5T6!", true},
		{"No numbers", "Qlowe_r_t_y!", true},
		{"No special characters", "Q1we3r4r5t6", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswordStrength(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePasswordStrength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
