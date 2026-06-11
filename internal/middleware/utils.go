package middleware

import "net/http"

// Преобразует middleware из формата func(http.HandlerFunc) http.HandlerFunc
// в стандартный формат chi/net/http: func(http.Handler) http.Handler.
func ToMiddleware(fn func(http.HandlerFunc) http.HandlerFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// Оборачиваем вызов следующего хендлера в замыкание
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Адаптируем метод ServeHTTP следующего обработчика к типу http.HandlerFunc
			adaptedNext := http.HandlerFunc(next.ServeHTTP)

			// Вызываем наше кастомное middleware fn, передавая ему адаптированный хендлер.
			// Это возвращает новый http.HandlerFunc.
			resultHandlerFunc := fn(adaptedNext)

			// Выполняем полученный результирующий хендлер с текущими ResponseWriter и Request
			resultHandlerFunc(w, r)
		})
	}
}
