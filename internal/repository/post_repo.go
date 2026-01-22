package repository

import (
	"advanced-blog-management-system/internal/model"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type PostRepo struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{db: db}
}

// TODO: Реализовать Create()
// INSERT INTO posts (title, content, author_id, created_at) VALUES (...) RETURNING id
// Установить CreatedAt = time.Now(), сканировать ID в post
func (r *PostRepo) Create(ctx context.Context, post *model.Post) error {
	// TODO: реализовать
	return nil
}

// TODO: Реализовать GetByID()
// SELECT ... FROM posts WHERE id = $1
// Если не найдена - вернуть nil, nil; если ошибка - вернуть fmt.Errorf()
func (r *PostRepo) GetByID(ctx context.Context, id int) (*model.Post, error) {
	// TODO: реализовать
	return nil, nil
}

// TODO: Реализовать GetAll()
// SELECT ... FROM posts ORDER BY created_at DESC LIMIT $1 OFFSET $2
func (r *PostRepo) GetAll(ctx context.Context, limit, offset int) ([]*model.Post, error) {
	// TODO: реализовать
	return nil, nil
}

// TODO: Реализовать GetTotalCount()
// SELECT COUNT(*) FROM posts
func (r *PostRepo) GetTotalCount(ctx context.Context) (int, error) {
	// TODO: реализовать
	return 0, nil
}

// TODO: Реализовать Exists()
// SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)
func (r *PostRepo) Exists(ctx context.Context, id int) (bool, error) {
	// TODO: реализовать
	return false, nil
}

// TODO: Реализовать GetByAuthorID()
// SELECT ... FROM posts WHERE author_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
func (r *PostRepo) GetByAuthorID(ctx context.Context, authorID int, limit, offset int) ([]*model.Post, error) {
	// TODO: реализовать
	return nil, nil
}

// TODO: Реализовать GetTotalCountByAuthorID()
// SELECT COUNT(*) FROM posts WHERE author_id = $1
func (r *PostRepo) GetTotalCountByAuthorID(ctx context.Context, authorID int) (int, error) {
	// TODO: реализовать
	return 0, nil
}
