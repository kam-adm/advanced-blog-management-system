package handler

import (
	"advanced-blog-management-system/internal/errors/apperrors"
	"advanced-blog-management-system/internal/logger"
	"advanced-blog-management-system/internal/middleware"
	"advanced-blog-management-system/internal/model"
	"advanced-blog-management-system/internal/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type CommentHandler struct {
	commentService service.CommentServiceInterface
	eventLogger    *logger.EventLogger
}

func NewCommentHandler(commentService service.CommentServiceInterface, eventLogger *logger.EventLogger) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		eventLogger:    eventLogger,
	}
}

// Create обрабатывает POST /api/posts/{postId}/comments
// TODO: Получить userID из контекста, извлечь postID, распарсить JSON, валидировать,
// вызвать commentService.Create(). Вернуть CommentResponse (201) или ошибку
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать
}

// GetByPost обрабатывает GET /api/posts/{postId}/comments?limit=10&offset=0
// TODO: Извлечь postID, получить limit и offset, валидировать,
// вызвать commentService.GetByPost(). Вернуть список комментариев и общее количество
func (h *CommentHandler) GetByPost(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать
}

// respondWithJSON отправляет JSON ответ
// TODO: Установить Content-Type, WriteHeader, закодировать JSON
func (h *CommentHandler) respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	// TODO: реализовать
}

// respondWithError отправляет JSON ошибку
// TODO: Создать ErrorResponse и отправить используя respondWithJSON()
func (h *CommentHandler) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	// TODO: реализовать
}
