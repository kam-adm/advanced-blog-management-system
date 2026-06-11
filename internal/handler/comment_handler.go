package handler

import (
	"advanced-blog-management-system/internal/errors/apperrors"
	"advanced-blog-management-system/internal/logger"
	"advanced-blog-management-system/internal/middleware"
	"advanced-blog-management-system/internal/model"
	"advanced-blog-management-system/internal/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CommentHandler struct {
	commentService service.CommentServiceInterface
	eventLogger    *logger.EventLogger
}

// Создает новый экземпляр CommentHandler
func NewCommentHandler(commentService service.CommentServiceInterface, eventLogger *logger.EventLogger) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		eventLogger:    eventLogger,
	}
}

// Обрабатывает POST запросы на создание комментариев
func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	// 1. Извлекаем userID из middleware
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok || userID == 0 {
		WriteError(w, apperrors.ErrUnauthorized.Error(), http.StatusUnauthorized)
		return
	}

	// 2. Извлекаем и валидируем postID
	postIDStr := chi.URLParam(r, "postId")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		WriteError(w, apperrors.ErrInvalidPostID.Error(), http.StatusBadRequest)
		return
	}

	// 3. Распарсиваем тело запроса
	var input model.CommentCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Инъекция post_id в структуру для работы валидатора
	input.PostID = postID

	// 4. Вызываем встроенную валидацию
	if err := input.Validate(); err != nil {
		HandleServiceError(w, err)
		return
	}

	// 5. Вызываем сервис в строгом соответствии с его сигнатурой
	comment, err := h.commentService.Create(r.Context(), userID, postID, input.Content)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	// 6. Отправляем асинхронное событие в канал логгера
	if h.eventLogger != nil {
		logMessage := fmt.Sprintf("user %d created comment %d", userID, comment.ID)
		h.eventLogger.LogEvent(logMessage)
	}

	// Отправляем успешный ответ клиенту
	h.respondWithJSON(w, comment, http.StatusCreated)
}

// Обрабатывает GET запросы на получение списка комментариев
func (h *CommentHandler) GetByPost(w http.ResponseWriter, r *http.Request) {
	// Извлекаем и валидируем postID из URL
	postIDStr := chi.URLParam(r, "postId")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		WriteError(w, apperrors.ErrInvalidPostID.Error(), http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// 3. Вызываем сервис получения данных
	comments, total, err := h.commentService.GetByPost(r.Context(), postID, limit, offset)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	// Формируем успешный ответ со списком и метаданными
	response := map[string]interface{}{
		"comments": comments,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	}

	h.respondWithJSON(w, response, http.StatusOK)
}

// Отправляет JSON ответ клиенту
func (h *CommentHandler) respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
