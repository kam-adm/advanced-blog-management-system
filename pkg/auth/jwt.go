package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

// Представляет данные, хранимые в JWT токене
type Claims struct {
	UserID   int    `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Управляет полным жизненным циклом JWT токенов
type JWTManager struct {
	secretKey []byte
	ttl       time.Duration
}

// Инициализирует менеджер токенов с секретным ключом и TTL в часах
func NewJWTManager(secretKey string, ttlHours int) *JWTManager {
	return &JWTManager{
		secretKey: []byte(secretKey),
		ttl:       time.Duration(ttlHours) * time.Hour,
	}
}

// Создает новый подписанный токен, возвращает строку, время истечения и ошибку
func (m *JWTManager) GenerateToken(userID int, email, username string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(m.ttl)

	claims := Claims{
		UserID:   userID,
		Email:    email,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "advanced-blog-api",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	// Создаем токен с использованием стандартного алгоритма HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен секретным ключом
	tokenString, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// Парсит токен, проверяет подпись и срок его действия
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	// Парсим токен и сразу проверяем метод подписи (HMAC)
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.secretKey, nil
	})

	// Обработка ошибок парсинга (включая истечение срока действия)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	// Извлекаем Claims и проверяем валидность токена
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid || claims == nil {
		return nil, ErrInvalidToken
	}

	// Дополнительная явная проверка времени действия на всякий случай
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, ErrExpiredToken
	}

	return claims, nil
}
