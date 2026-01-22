package repository

import (
	"advanced-blog-management-system/internal/model"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type CommentRepo struct {
	db *sql.DB
}

func NewCommentRepo(db *sql.DB) *CommentRepo {
	return &CommentRepo{db: db}
}

// TODO: Реализовать Create()
// INSERT INTO comments (content, post_id, author_id, created_at, updated_at) VALUES (...) RETURNING id
// Установить CreatedAt, UpdatedAt = time.Now(), сканировать ID в comment
func (r *CommentRepo) Create(ctx context.Context, comment *model.Comment) error {
	// TODO: реализовать
	return nil
}

// TODO: Реализовать GetByID()
// SELECT ... FROM comments WHERE id = $1
// Если не найден - вернуть nil, nil; если ошибка - вернуть fmt.Errorf()
func (r *CommentRepo) GetByID(ctx context.Context, id int) (*model.Comment, error) {
	// TODO: реализовать
	return nil, nil
}

// TODO: Реализовать GetByPostID()
// SELECT ... FROM comments WHERE post_id = $1
// ORDER BY created_at ASC LIMIT $2 OFFSET $3
func (r *CommentRepo) GetByPostID(ctx context.Context, postID int, limit, offset int) ([]*model.Comment, error) {
	// TODO: реализовать
	return nil, nil
}

// TODO: Реализовать GetCountByPostID()
// SELECT COUNT(*) FROM comments WHERE post_id = $1
func (r *CommentRepo) GetCountByPostID(ctx context.Context, postID int) (int, error) {
	// TODO: реализовать
	return 0, nil
}
