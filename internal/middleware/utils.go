package middleware

import "net/http"

// TODO: Реализовать ToMiddleware()
// Функция преобразует middleware функцию из одного формата в другой.
// Нужна для совместимости разных типов middleware в chi.Router
//
// Сигнатура входного middleware:
//   func(http.HandlerFunc) http.HandlerFunc
//
// Сигнатура выходного middleware:
//   func(http.Handler) http.Handler
//
// Функция должна:
// 1. Принять middleware функцию fn типа func(http.HandlerFunc) http.HandlerFunc
// 2. Вернуть функцию типа func(http.Handler) http.Handler
// 3. Внутри возвращаемой функции:
//    - Принять handler типа http.Handler
//    - Вернуть новый http.Handler который:
//      * Вызывает fn(handler.ServeHTTP)
//      * Это возвращает http.HandlerFunc
//      * Вызывает эту функцию с (w, r)
//
// Пример использования в main.go:
//   router.Use(middleware.ToMiddleware(authMiddleware.RequireAuth))
//
// Это позволяет использовать middleware типа RequireAuth(http.HandlerFunc) http.HandlerFunc
// с chi.Router который требует middleware типа func(http.Handler) http.Handler
func ToMiddleware(fn func(http.HandlerFunc) http.HandlerFunc) func(http.Handler) http.Handler {
	// TODO: реализовать
	return nil
}
