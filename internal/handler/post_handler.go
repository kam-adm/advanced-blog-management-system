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

type PostHandler struct {
	postService service.PostServiceInterface
	eventLogger *logger.EventLogger
}

// Создает новый экземпляр PostHandler
func NewPostHandler(postService service.PostServiceInterface, eventLogger *logger.EventLogger) *PostHandler {
	return &PostHandler{
		postService: postService,
		eventLogger: eventLogger,
	}
}

// Обрабатывает POST
func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	// 1. Извлекаем ID авторизованного пользователя из middleware
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok || userID == 0 {
		WriteError(w, apperrors.ErrUnauthorized.Error(), http.StatusUnauthorized)
		return
	}

	// 2. Декодируем тело запроса в правильный Request
	var input model.PostCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 3. Вызываем встроенную валидацию
	if err := input.Validate(); err != nil {
		HandleServiceError(w, err)
		return
	}

	// 4. Передаем на уровень бизнес-логики в строгом соответствии с сигнатурой метода
	post, err := h.postService.Create(r.Context(), userID, &input)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	// 5. Отправляем асинхронное событие в канал через LogEvent
	if h.eventLogger != nil {
		logMessage := fmt.Sprintf("user %d created post %d", userID, post.ID)
		h.eventLogger.LogEvent(logMessage)
	}

	h.respondWithJSON(w, post, http.StatusCreated)
}

// Обрабатывает GET
func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// 1. Парсим ID поста из параметров роутера
	postIDStr := chi.URLParam(r, "id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		WriteError(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// 2. Получаем userID из контекста, если пользователь авторизован
	userID, _ := middleware.GetUserIDFromContext(r.Context())

	// 3. Запрашиваем пост через сервис
	post, err := h.postService.GetByID(r.Context(), postID, userID)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	h.respondWithJSON(w, post, http.StatusOK)
}

// Обрабатывает GET
func (h *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// 1. Парсим параметры и задаем дефолтные значения
	limit := 10
	if lStr := r.URL.Query().Get("limit"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
			limit = l
			if limit > 100 { // Ограничиваем максимальный размер пачки в 100 элементов
				limit = 100
			}
		}
	}

	offset := 0
	if oStr := r.URL.Query().Get("offset"); oStr != "" {
		if o, err := strconv.Atoi(oStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// 2. Вызываем сервис для выборки списка
	posts, total, err := h.postService.GetAll(r.Context(), limit, offset)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	response := map[string]interface{}{
		"posts":  posts,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}

	h.respondWithJSON(w, response, http.StatusOK)
}

// Обрабатывает GET
func (h *PostHandler) GetByAuthor(w http.ResponseWriter, r *http.Request) {
	// 1. Извлекаем ID автора из URL
	authorIDStr := chi.URLParam(r, "authorID")
	authorID, err := strconv.Atoi(authorIDStr)
	if err != nil || authorID <= 0 {
		WriteError(w, "Invalid author ID", http.StatusBadRequest)
		return
	}

	// 2. Парсим параметры
	limit := 10
	if lStr := r.URL.Query().Get("limit"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := 0
	if oStr := r.URL.Query().Get("offset"); oStr != "" {
		if o, err := strconv.Atoi(oStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// 3. Запрос к бизнес-слою
	posts, total, err := h.postService.GetByAuthor(r.Context(), authorID, limit, offset)
	if err != nil {
		HandleServiceError(w, err)
		return
	}

	response := map[string]interface{}{
		"posts":  posts,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}

	h.respondWithJSON(w, response, http.StatusOK)
}

// Отправляет JSON ответ
func (h *PostHandler) respondWithJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
