package handler

import (
	"advanced-blog-management-system/internal/errors/apperrors"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// Представляет структуру ошибки в ответе
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// Отправляет стандартизированный JSON-ответ об ошибке
func WriteError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	}

	// Кодируем напрямую в http.ResponseWriter для экономии памяти
	_ = json.NewEncoder(w).Encode(resp)
}

// Преобразует бизнес-ошибки приложения в корректные HTTP-коды и ответы
func HandleServiceError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	// 1. Проверяем ошибки валидации от библиотеки go-playground/validator
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		msg := "Validation failed: "
		for _, vErr := range validationErrors {
			msg += fmt.Sprintf("[%s: failed '%s' validation] ", vErr.Field(), vErr.Tag())
		}
		WriteError(w, msg, http.StatusBadRequest)
		return
	}

	// 2. Сопоставляем внутренние ошибки apperrors с HTTP статусами через errors.Is()
	switch {
	case errors.Is(err, apperrors.ErrUserAlreadyExists):
		WriteError(w, err.Error(), http.StatusConflict)

	case errors.Is(err, apperrors.ErrInvalidCredentials) || errors.Is(err, apperrors.ErrUnauthorized):
		WriteError(w, err.Error(), http.StatusUnauthorized)

	case errors.Is(err, apperrors.ErrPostNotFound) || errors.Is(err, apperrors.ErrCommentNotFound):
		WriteError(w, err.Error(), http.StatusNotFound)

	case errors.Is(err, apperrors.ErrForbidden):
		WriteError(w, err.Error(), http.StatusForbidden)

	case errors.Is(err, apperrors.ErrInvalidPostID):
		WriteError(w, err.Error(), http.StatusBadRequest)

	// 3. Все остальные непредвиденные системные ошибки маскируем под 500 Internal Server Error
	default:
		WriteError(w, "An internal server error occurred", http.StatusInternalServerError)
	}
}
