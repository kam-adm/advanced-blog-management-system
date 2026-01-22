package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

// Claims представляет данные, хранимые в JWT токене
type Claims struct {
	UserID   int    `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// TODO: Добавить структуру JWTManager с полями secretKey []byte и ttl time.Duration
// Реализовать конструктор NewJWTManager(secretKey string, ttlHours int) *JWTManager
type JWTManager struct {
	secretKey []byte
	ttl       time.Duration
}

// TODO: Реализовать конструктор NewJWTManager
// Параметры: secretKey string, ttlHours int
// Вернуть: *JWTManager с инициализированными полями
func NewJWTManager(secretKey string, ttlHours int) *JWTManager {
	// TODO: реализовать
	return nil
}

// TODO: Реализовать метод GenerateToken(userID int, email, username string) (string, time.Time, error)
// - Создать Claims со сроком действия = now + m.ttl
// - Установить IssuedAt, NotBefore, ExpiresAt, Issuer, Subject в Claims
// - Подписать токен методом HS256 секретным ключом
// - Вернуть (tokenString, expiresAt, nil) или ("", time.Time{}, error)
func (m *JWTManager) GenerateToken(userID int, email, username string) (string, time.Time, error) {
	// TODO: реализовать
	return "", time.Time{}, nil
}

// TODO: Реализовать метод ValidateToken(tokenString string) (*Claims, error)
// - Распарсить токен с проверкой подписи через jwt.ParseWithClaims
// - Извлечь Claims из токена
// - Проверить что token.Valid == true и claims != nil
// - Проверить что время истечения еще не наступило (ExpiresAt.Time > now)
// - Вернуть claims или ошибку (ErrInvalidToken, ErrExpiredToken)
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	// TODO: реализовать
	return nil, nil
}
