package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"advanced-blog-management-system/internal/handler"
	"advanced-blog-management-system/internal/logger"
	"advanced-blog-management-system/internal/middleware"
	"advanced-blog-management-system/internal/repository"
	"advanced-blog-management-system/internal/service"
	"advanced-blog-management-system/pkg/auth"
	"advanced-blog-management-system/pkg/database"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

// Хранит все настройки приложения
type Config struct {
	DBConfig   database.Config
	JWTSecret  string
	ServerAddr string
}

func main() {
	// 1. Загрузка переменных окружения из .env файла
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}
	log.Println("Configuration loaded successfully")

	// 2. Чтение конфигурации из переменных окружения
	cfg := loadConfig()
	if cfg.JWTSecret == "" {
		log.Fatal("Fatal: JWT_SECRET environment variable is required")
	}

	// 3. Создание подключения к PostgreSQL
	log.Println("Connecting to PostgreSQL...")
	db, err := database.NewPostgresDB(cfg.DBConfig)
	if err != nil {
		log.Fatalf("Fatal: database connection failed: %v", err)
	}
	defer database.Close(db)

	// 4. Выполнение миграций БД
	log.Println("Running database migrations...")
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Fatal: database migrations failed: %v", err)
	}

	// Инициализация менеджера JWT токенов
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, 24)

	// Инициализация асинхронного фонового логгера событий на горутинах и каналах
	eventLogger := logger.NewEventLogger("data/log.txt")
	eventLogger.Start()

	// 5. Инициализация репозиториев
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepo(db)
	commentRepo := repository.NewCommentRepo(db)

	// 6. Инициализация сервисов
	userService := service.NewUserService(userRepo, jwtManager)
	postService := service.NewPostService(postRepo, userRepo)
	commentService := service.NewCommentService(commentRepo, postRepo)

	// 7. Инициализация обработчиков
	authHandler := handler.NewAuthHandler(userService)
	postHandler := handler.NewPostHandler(postService, eventLogger)
	commentHandler := handler.NewCommentHandler(commentService, eventLogger)

	// 8. Создать и настроить HTTP роутер
	r := setupRouter(authHandler, postHandler, commentHandler, jwtManager)

	// 9. Создать HTTP сервер
	srv := &http.Server{
		Addr:         cfg.ServerAddr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 10. Запустить сервер в отдельной горутине
	go func() {
		log.Printf("Server is running on http://%s", cfg.ServerAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Fatal: server failed to start: %v", err)
		}
	}()

	// 11. Установить обработчик сигналов завершения (SIGINT, SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Ожидание сигнала
	sig := <-quit
	log.Printf("Received signal: %v. Starting graceful shutdown...", sig)

	// 12. Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Останавливаем HTTP сервер
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Fatal: server forced to shutdown: %v", err)
	}

	// Корректно завершаем работу логгера
	eventLogger.Stop()

	log.Println("Server stopped gracefully. Goodbye!")
}

// Вспомогательная функция для сборки конфигурации с безопасной конвертацией типов
func loadConfig() Config {
	getEnv := func(key, fallback string) string {
		if value, exists := os.LookupEnv(key); exists {
			return value
		}
		return fallback
	}

	// Конвертируем порт в число int под требования структуры database.Config
	portStr := getEnv("DB_PORT", "5432")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 5432
	}

	dbCfg := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     port,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "blog_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	serverHost := getEnv("SERVER_HOST", "0.0.0.0")
	serverPort := getEnv("SERVER_PORT", "8080")

	return Config{
		DBConfig:   dbCfg,
		JWTSecret:  os.Getenv("JWT_SECRET"),
		ServerAddr: fmt.Sprintf("%s:%s", serverHost, serverPort),
	}
}

// Настройка маршрутизатора и привязка Middleware
func setupRouter(
	authHandler *handler.AuthHandler,
	postHandler *handler.PostHandler,
	commentHandler *handler.CommentHandler,
	jwtManager *auth.JWTManager,
) *chi.Mux {
	r := chi.NewRouter()

	// Инициализируем системные логирующие middleware
	sysLogger := log.New(os.Stdout, "", log.LstdFlags)
	logMW := middleware.NewLoggingMiddleware(sysLogger)
	authMW := middleware.NewAuthMiddleware(jwtManager)

	// Регистрация глобальных Middleware
	r.Use(logMW.Logger)
	r.Use(logMW.Recovery)
	r.Use(logMW.CORS)

	// Публичные эндпоинты API
	r.Get("/api/health", handler.HealthCheckHandler)
	r.Post("/api/register", authHandler.Register)
	r.Post("/api/login", authHandler.Login)

	r.Get("/api/posts", postHandler.GetAll)
	r.Get("/api/posts/{id}", postHandler.GetByID)
	r.Get("/api/users/{authorID}/posts", postHandler.GetByAuthor)
	r.Get("/api/posts/{postId}/comments", commentHandler.GetByPost)

	// Защищенные эндпоинты
	r.Group(func(r chi.Router) {
		r.Use(middleware.ToMiddleware(authMW.RequireAuth))

		// Управление публикациями постов
		r.Post("/api/posts", postHandler.Create)

		// Управление комментариями
		r.Post("/api/posts/{postId}/comments", commentHandler.Create)
	})

	return r
}
