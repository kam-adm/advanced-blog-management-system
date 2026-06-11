package repository

import (
	"advanced-blog-management-system/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type CommentRepo struct {
	db *sql.DB
}

// Создает новый экземпляр репозитория комментариев
func NewCommentRepo(db *sql.DB) *CommentRepo {
	return &CommentRepo{db: db}
}

// Сохраняет новый комментарий в базу данных
func (r *CommentRepo) Create(ctx context.Context, comment *model.Comment) error {
	query := `
		INSERT INTO comments (content, post_id, author_id, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id`

	now := time.Now()
	comment.CreatedAt = now
	comment.UpdatedAt = now

	// Используем QueryRowContext для безопасного параметризованного запроса и сканирования RETURNING id
	err := r.db.QueryRowContext(ctx, query, comment.Content, comment.PostID, comment.AuthorID, comment.CreatedAt, comment.UpdatedAt).Scan(&comment.ID)
	if err != nil {
		return fmt.Errorf("failed to insert comment: %w", err)
	}

	return nil
}

// Находит и возвращает комментарий по его уникальному идентификатору
func (r *CommentRepo) GetByID(ctx context.Context, id int) (*model.Comment, error) {
	query := `SELECT id, content, post_id, author_id, created_at, updated_at FROM comments WHERE id = $1`

	var comment model.Comment
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&comment.ID,
		&comment.Content,
		&comment.PostID,
		&comment.AuthorID,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Согласно ТЗ, если сущность не найдена — возвращаем nil без ошибки
		}
		return nil, fmt.Errorf("failed to query comment by id: %w", err)
	}

	return &comment, nil
}

// Выбирает список комментариев к посту с поддержкой постраничной навигации
func (r *CommentRepo) GetByPostID(ctx context.Context, postID int, limit, offset int) ([]*model.Comment, error) {
	query := `
		SELECT id, content, post_id, author_id, created_at, updated_at 
		FROM comments 
		WHERE post_id = $1 
		ORDER BY created_at ASC 
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, postID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments by post_id: %w", err)
	}
	defer rows.Close()

	comments := make([]*model.Comment, 0, limit)
	for rows.Next() {
		var comment model.Comment
		err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.PostID,
			&comment.AuthorID,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment row: %w", err)
		}
		comments = append(comments, &comment)
	}

	// Проверяем, не возникло ли ошибок во время итерации по строкам
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during comments rows iteration: %w", err)
	}

	return comments, nil
}

// Подсчитывает общее число комментариев к конкретному посту (необходимо для метаданных пагинации)
func (r *CommentRepo) GetCountByPostID(ctx context.Context, postID int) (int, error) {
	query := `SELECT COUNT(*) FROM comments WHERE post_id = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, postID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count comments by post_id: %w", err)
	}

	return count, nil
}
