# Система управления блогом (Advanced Blog Management System)

Веб-приложение на языке программирования Go версии 1.24.5 для управления блог-платформой.

## Технологический стек проекта
- **Язык разработки:** Go 1.24.5
- **Маршрутизация и HTTP-инфраструктура:** Chi Router (v5)
- **Система управления базами данных:** PostgreSQL 15 (Alpine)
- **Аутентификация и сессии:** JWT
- **Криптографическое хеширование:** bcrypt
- **Валидация входящих данных:** go-playground/validator
- **Контейнеризация и оркестрация:** Docker, Docker Compose
- **Управление конфигурацией среды:** godotenv (.env)

## Архитектурные слои и структура компонентов
Приложение разделено на изолированные слои с четкими зонами ответственности:
1. **cmd/api/main.go** — точка входа в приложение. Отвечает за сборку конфигурации, инициализацию пула соединений СУБД, запуск фоновых воркеров и HTTP-сервера с механизмами корректного завершения работы (Graceful Shutdown).
2. **internal/handler/** — слой представления. Обеспечивает десериализацию JSON-тел запросов, извлечение параметров маршрутизации, запуск встроенной валидации DTO и вызов сервисов бизнес-логики.
3. **internal/service/** — слой бизнес-логики. Содержит правила проверки уникальности сущностей, валидацию криптографической стойкости паролей и алгоритмы выпуска токенов сессий. Взаимодействует со слоем данных исключительно через абстракции (интерфейсы).
4. **internal/repository/** — слой доступа к данным. Реализует прямые параметризованные SQL-запросы к СУБД PostgreSQL с контролем таймаутов через context.Context.
5. **internal/middleware/** — сервисные сетевые компоненты. Реализуют установку CORS-заголовков, структурированное логирование HTTP-трафика, авторизацию сессий и перехват критических паник приложения (Recovery).
6. **internal/logger/** — компонент асинхронного логирования действий пользователей. Осуществляет неблокирующее чтение системных событий из буферизованного канала и их запись в файл data/log.txt в отдельной горутине.

---

## Инструкция по развертыванию приложения

1. **Создание файла конфигурации среды**
   Скопируйте демонстрационные переменные окружения в рабочий файл настройки:
   ```bash
   cp .env.example .env
   ```

2. **Запуск инфраструктуры в контейнерах**
   Выполните сборку образов и запуск контейнеров в фоновом режиме. Исполняемый Go-модуль начнет работу автоматически после успешного прохождения проверки доступности СУБД (healthcheck):
   ```bash
   docker-compose up --build -d
   ```

3. **Мониторинг системных логов**
   ```bash
   docker-compose logs -f app
   ```

4. **Остановка сервисов и удаление томов данных**
   ```bash
   docker-compose down -v
   ```

---

## Сценарий верификации и ручного тестирования API

### 1. Проверка доступности сервиса (Health Check)
```bash
curl -X GET http://localhost:8080/api/health
```
*Ожидаемый ответ:* `{"status":"ok"}`

### 2. Регистрация нового пользователя в системе
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user",
    "email": "user@mail.ru",
    "password": "Q1we3r4r5t6!"
  }'
```

### 3. Аутентификация пользователя и получение токена сессии
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@mail.ru",
    "password": "Q1we3r4r5t6!"
  }'
```
*Значение строки "token" из ответа СУБД необходимо использовать в следующих запросах вместо подстановки `ВСТАВЬТЕ_ТОКЕН`.*

### 4. Создание новой публикации с текстом благодарности преподавателям
```bash
curl -X POST http://localhost:8080/api/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ВСТАВЬТЕ_ТОКЕН" \
  -d '{
    "title": "Слова искренней благодарности преподавательскому составу",
    "content": "Выражаю глубокую признательность всем преподавателям за их высокий профессионализм, терпение и чуткое руководство на протяжении всего процесса обучения. Благодаря вашему труду я получил прочные фундаментальные знания, освоил современные технологии и успешно разработал данную выпускную квалификационную работу!"
  }'
```

### 5. Получение списка всех публикаций с пагинацией
```bash
curl -X GET "http://localhost:8080/api/posts?limit=5&offset=0"
```

### 6. Добавление комментария к публикации
```bash
curl -X POST http://localhost:8080/api/posts/1/comments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ВСТАВЬТЕ_ТОКЕН" \
  -d '{
    "content": "This is an amazing diploma project base structure!"
  }'
```

### 7. Чтение списка комментариев к публикации
```bash
curl -X GET "http://localhost:8080/api/posts/1/comments?limit=10&offset=0"
```

---

## Модульное тестирование (Unit-тестирование)

Для проверки работоспособности бизнес-логики приложения в изоляции от инфраструктуры баз данных реализовано модульное тестирование критических путей выполнения с применением Mock-заглушек.

### Команда для запуска всех тест-кейсов проекта
```bash
go test ./... -v
```

### Отчет о прохождении автоматизированного тестирования
```text
=== RUN   TestJWTManager_Lifecycle
--- PASS: TestJWTManager_Lifecycle (0.00s)
=== RUN   TestPassword_HashingAndVerification
--- PASS: TestPassword_HashingAndVerification (0.01s)
=== RUN   TestPassword_StrengthValidation
--- PASS: TestPassword_StrengthValidation (0.00s)
=== RUN   TestUserCreateRequest_Validation
--- PASS: TestUserCreateRequest_Validation (0.00s)
=== RUN   TestAuthMiddleware_RequireAuth
--- PASS: TestAuthMiddleware_RequireAuth (0.00s)
=== RUN   TestPostService_Create_Success
--- PASS: TestPostService_Create_Success (0.00s)
PASS
ok      advanced-blog-management-system/pkg/auth               0.012s
ok      advanced-blog-management-system/internal/model         0.003s
ok      advanced-blog-management-system/internal/middleware    0.005s
ok      advanced-blog-management-system/internal/service       0.004s
```
