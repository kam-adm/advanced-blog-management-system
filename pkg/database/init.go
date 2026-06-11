package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// Config содержит конфигурацию базы данных
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresDB создает, проверяет и настраивает пул подключений к PostgreSQL
func NewPostgresDB(cfg Config) (*sql.DB, error) {
	// 1. Формируем DSN строку
	dsn := GetDSN(cfg)

	// 2. Открываем подключение к драйверу
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// 3. Реально проверяем сетевую доступность СУБД через Ping
	if err := db.Ping(); err != nil {
		_ = db.Close() // Закрываем дескриптор в случае неудачи
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// 4. Настраиваем пул соединений для высокой производительности и защиты от утечек
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// 5. Логируем успешное событие
	log.Println("Connected to PostgreSQL database")

	// 6. Возвращаем рабочий пул соединений
	return db, nil
}

// GetDSN формирует Data Source Name строку для подключения к PostgreSQL
func GetDSN(cfg Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)
}

// Close безопасно закрывает соединение с базой данных
func Close(db *sql.DB) {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("Error while closing database connection: %v", err)
			return
		}
		log.Println("Database connection closed")
	}
}

// TestConnection выполняет тестовый холостой запрос к БД для проверки стабильности подключения
func TestConnection(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database instance is nil")
	}

	// Выполняем самый легкий и быстрый SQL запрос для проверки
	if err := db.Ping(); err != nil {
		return fmt.Errorf("database connection test failed: %w", err)
	}

	return nil
}
