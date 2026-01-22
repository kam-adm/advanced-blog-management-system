package apperrors

import "errors"

// TODO: Определить переменные ошибок приложения
// Используйте errors.New() для создания каждой ошибки.
// Необходимы следующие ошибки:
// - ErrUnauthorized - пользователь не аутентифицирован
// - ErrForbidden - пользователь не имеет прав доступа
// - ErrUserNotFound - пользователь не найден в БД
// - ErrUserAlreadyExists - пользователь с таким email/username уже существует
// - ErrInvalidCredentials - неверный email или пароль при входе
// - ErrPostNotFound - пост не найден в БД
// - ErrCommentNotFound - комментарий не найден в БД
// - ErrInvalidPostID - передан невалидный ID поста
//
// Эти переменные будут использоваться в:
// - Services: возвращать эти ошибки для обозначения бизнес-логики
// - Handlers: проверять ошибки через errors.Is() и возвращать нужные HTTP коды
//
// Пример создания:
// var (
//     ErrExample = errors.New("example error message")
//     ErrOther = errors.New("other error message")
// )
