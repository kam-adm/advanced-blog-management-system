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

type PostHandler struct {
	postService service.PostServiceInterface
	eventLogger *logger.EventLogger
}

// NewPostHandler создает новый экземпляр PostHandler
func NewPostHandler(postService service.PostServiceInterface, eventLogger *logger.EventLogger) *PostHandler {
	return &PostHandler{
		postService: postService,
		eventLogger: eventLogger,
	}
}

// Create обрабатывает POST /api/posts
// TODO: Получить userID из контекста, распарсить JSON, валидировать, вызвать postService.Create()
// Вернуть PostResponse со статусом 201 или ошибку
func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать
}

// GetByID обрабатывает GET /api/posts/{id}
// TODO: Извлечь postID из URL, получить userID из контекста (может быть 0),
// вызвать postService.GetByID(). Вернуть PostResponse (200) или ошибку
func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать
}

// GetAll обрабатывает GET /api/posts?limit=10&offset=0
// TODO: Получить limit и offset из query, валидировать (limit 1-100, default 10),
// вызвать postService.GetAll(). Вернуть список постов и общее количество
func (h *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать
}

// GetByAuthor обрабатывает GET /api/posts/author/{authorID}?limit=10&offset=0
// TODO: Извлечь authorID из URL, получить limit и offset из query,
// вызвать postService.GetByAuthor(). Вернуть список постов автора и общее количество
func (h *PostHandler) GetByAuthor(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать
}

// respondWithJSON отправляет JSON ответ
// TODO: Установить Content-Type, WriteHeader, закодировать JSON
func (h *PostHandler) respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	// TODO: реализовать
}

// respondWithError отправляет JSON ошибку
// TODO: Создать ErrorResponse и отправить используя respondWithJSON()
func (h *PostHandler) respondWithError(w http.ResponseWriter, message string, statusCode int) {
	// TODO: реализовать
}
