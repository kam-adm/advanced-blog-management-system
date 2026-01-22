package handler

import (
	"advanced-blog-management-system/internal/errors/apperrors"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// ErrorResponse представляет структуру ошибки в ответе
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// TODO: Реализовать WriteError(w http.ResponseWriter, message string, statusCode int)
// 1. Установить Content-Type: application/json
// 2. Установить статус код через WriteHeader()
// 3. Закодировать ErrorResponse в JSON используя json.NewEncoder()
// ErrorResponse должен содержать Error = http.StatusText(statusCode) и Message = message
func WriteError(w http.ResponseWriter, message string, statusCode int) {
	// TODO: реализовать
}

// TODO: Реализовать HandleServiceError(w http.ResponseWriter, err error)
// 1. Проверить тип ошибки используя errors.Is()
// 2. Для каждого типа ошибки из apperrors вернуть нужный HTTP статус:
//    - ErrUserAlreadyExists -> 409 Conflict
//    - ErrInvalidCredentials -> 401 Unauthorized
//    - ErrPostNotFound -> 404 Not Found
//    - ErrCommentNotFound -> 404 Not Found
//    - ErrForbidden -> 403 Forbidden
//    - ErrUnauthorized -> 401 Unauthorized
// 3. Проверить validator.ValidationErrors -> 400 Bad Request
// 4. Для других ошибок -> 500 Internal Server Error
// 5. Использовать WriteError() для отправки ответа с нужным сообщением
func HandleServiceError(w http.ResponseWriter, err error) {
	// TODO: реализовать
}
