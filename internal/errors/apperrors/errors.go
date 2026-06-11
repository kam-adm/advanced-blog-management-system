package apperrors

import "errors"

// Переменные ошибок приложения
var (
	// Ошибки аутентификации и авторизации
	ErrUnauthorized       = errors.New("user is not authenticated")
	ErrForbidden          = errors.New("user does not have permission to perform this action")
	ErrInvalidCredentials = errors.New("invalid email or password")

	// Ошибки User
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user with this email or username already exists")

	// Ошибки Post
	ErrPostNotFound  = errors.New("post not found")
	ErrInvalidPostID = errors.New("provided post ID is invalid or malformed")

	// Ошибки Comment
	ErrCommentNotFound = errors.New("comment not found")
)
