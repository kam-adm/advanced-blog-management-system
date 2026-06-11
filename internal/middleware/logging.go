package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

type LoggingMiddleware struct {
	logger *log.Logger
}

func NewLoggingMiddleware(logger *log.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{logger: logger}
}

// responseWriter обертка для ResponseWriter чтобы перехватить статус код
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Логирует каждый входящий HTTP-запрос: IP, метод, путь, статус-код и время выполнения
func (m *LoggingMiddleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Оборачиваем стандартный ResponseWriter. По умолчанию выставляем 200 OK,
		// так как если WriteHeader не будет вызван явно, Go отдаст статус 200.
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		// Считаем время выполнения запроса
		duration := time.Since(start)

		// Записываем структурированный лог запроса
		m.logger.Printf(
			"[HTTP] %s - %s %s | Status: %d | Duration: %v",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			rw.statusCode,
			duration,
		)
	})
}

// Перехватывает паники в приложении, логирует стек вызовов и предотвращает падение сервера
func (m *LoggingMiddleware) Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Логируем ошибку паники и трассировку стека для отладки
				m.logger.Printf("[PANIC RECOVERY] %v\nStack Trace:\n%s", err, debug.Stack())

				// Отдаем клиенту стандартизированный JSON с кодом 500
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)

				resp := map[string]string{
					"error":   http.StatusText(http.StatusInternalServerError),
					"message": "A critical internal server error occurred",
				}
				_ = json.NewEncoder(w).Encode(resp)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// Настраивает заголовки cross-origin ресурсов и обрабатывает предзапросы OPTIONS
func (m *LoggingMiddleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400") // Кэширование CORS предзапроса на 24 часа

		// Если это предзапрос (preflight), завершаем его успешно со статусом 204 No Content
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
