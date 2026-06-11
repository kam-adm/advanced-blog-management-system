package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Находит, проверяет историю и последовательно применяет новые SQL миграции
func Migrate(db *sql.DB) error {
	migrationsDir := "migrations"

	// 1. Создаем таблицу для учета миграций, если её ещё нет в БД
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to initialize migrations table: %w", err)
	}

	// 2. Получаем список всех элементов внутри директории migrations/
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// 3. Итерируемся по списку (они уже отсортированы по алфавиту)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()

		// Фильтруем файлы, выбирая строго с расширением .sql
		if !strings.HasSuffix(strings.ToLower(filename), ".sql") {
			continue
		}

		// Проверяем, не была ли эта миграция применена ранее
		applied, err := isMigrationApplied(db, filename)
		if err != nil {
			return fmt.Errorf("failed to check migration status for %s: %w", filename, err)
		}

		if applied {
			log.Printf("Migration %s already applied, skipping", filename)
			continue
		}

		log.Printf("Running migration: %s", filename)

		// Считываем контент файла
		filePath := filepath.Join(migrationsDir, filename)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		// Выполняем SQL-код в базе данных
		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to run migration %s: %w", filename, err)
		}

		// Фиксируем успешное выполнение миграции в таблице истории
		if err := recordMigration(db, filename); err != nil {
			return fmt.Errorf("failed to record migration history for %s: %w", filename, err)
		}

		log.Printf("Successfully applied migration: %s", filename)
	}

	log.Println("All database migrations checked and applied successfully")
	return nil
}

// Создает таблицу истории миграций schema_migrations
func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id SERIAL PRIMARY KEY,
			version VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`
	_, err := db.Exec(query)
	return err
}

// Проверяет, существует ли запись о файле в таблице истории
func isMigrationApplied(db *sql.DB, version string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1);`
	var exists bool
	err := db.QueryRow(query, version).Scan(&exists)
	return exists, err
}

// Вставляет запись о выполненном файле в таблицу истории
func recordMigration(db *sql.DB, version string) error {
	query := `INSERT INTO schema_migrations (version) VALUES ($1);`
	_, err := db.Exec(query, version)
	return err
}
