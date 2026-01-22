package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// TODO: Реализовать Migrate()
// Функция выполняет SQL миграции из папки migrations/:
// 1. Получить путь к папке migrations:
//    - Используйте os.ReadDir("migrations") для чтения содержимого
//    - Если ошибка - вернуть fmt.Errorf("failed to read migrations directory: %w", err)
// 2. Отфильтровать SQL файлы:
//    - Прочитать все файлы из папки migrations
//    - Выбрать только файлы с расширением .sql
//    - Отсортировать по имени (001_*, 002_*, 003_*, и т.д.)
// 3. Для каждого файла миграции по порядку:
//    - Залогировать "Running migration: <filename>"
//    - Прочитать содержимое файла используя os.ReadFile(filepath.Join("migrations", filename))
//    - Выполнить SQL используя db.Exec(string(content))
//    - Если ошибка - залогировать и вернуть fmt.Errorf("failed to run migration %s: %w", filename, err)
//    - Если успешно - залогировать "Successfully applied migration: <filename>"
// 4. После всех миграций вернуть nil
//
// Примечание: порядок выполнения ВАЖЕН!
// Файлы должны выполняться в алфавитном порядке:
// - 001_init_schema.sql (создание таблиц)
// - 002_add_foreign_keys.sql (внешние ключи)
// - 003_create_indexes.sql (индексы)
func Migrate(db *sql.DB) error {
	// TODO: реализовать
	return nil
}
