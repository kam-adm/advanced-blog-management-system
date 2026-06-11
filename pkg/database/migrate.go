package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Migrate находит, считывает и последовательно применяет все SQL миграции из папки migrations/
func Migrate(db *sql.DB) error {
	migrationsDir := "migrations"

	// 1. Получаем список всех элементов внутри директории migrations/
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// 2. Итерируемся по списку (os.ReadDir возвращает элементы, уже отсортированные по алфавиту)
	for _, entry := range entries {
		// Игнорируем вложенные папки, если они есть
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()

		// Фильтруем файлы, выбирая строго с расширением .sql
		if !strings.HasSuffix(strings.ToLower(filename), ".sql") {
			continue
		}

		// 3. Логируем начало выполнения конкретной миграции
		log.Printf("Running migration: %s", filename)

		// Конструируем полный путь к файлу и считываем его бинарное содержимое
		filePath := filepath.Join(migrationsDir, filename)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", filename, err)
		}

		// Выполняем SQL-код в базе данных
		_, err = db.Exec(string(content))
		if err != nil {
			log.Printf("Error during execution of migration %s: %v", filename, err)
			return fmt.Errorf("failed to run migration %s: %w", filename, err)
		}

		// Логируем успешный накат скрипта
		log.Printf("Successfully applied migration: %s", filename)
	}

	// 4. После успешного наката всех найденных файлов возвращаем nil
	log.Println("All database migrations applied successfully")
	return nil
}
