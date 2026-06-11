package repository

import (
	"advanced-blog-management-system/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type PostRepo struct {
	db *sql.DB
}

// Создает новый экземпляр репозитория постов
func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{db: db}
}

// Создает новый пост в базе данных
func (r *PostRepo) Create(ctx context.Context, post *model.Post) error {
	query := `
		INSERT INTO posts (title, content, author_id, created_at) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id`

	post.CreatedAt = time.Now()

	err := r.db.QueryRowContext(ctx, query, post.Title, post.Content, post.AuthorID, post.CreatedAt).Scan(&post.ID)
	if err != nil {
		return fmt.Errorf("failed to insert post: %w", err)
	}

	return nil
}

// Возвращает пост по его ID
func (r *PostRepo) GetByID(ctx context.Context, id int) (*model.Post, error) {
	query := `SELECT id, title, content, author_id, created_at FROM posts WHERE id = $1`

	var post model.Post
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.AuthorID,
		&post.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Объект не найден — по контракту возвращаем чистые (nil, nil)
		}
		return nil, fmt.Errorf("failed to query post by id: %w", err)
	}

	return &post, nil
}

// Получает все посты из базы данных с поддержкой пагинации
func (r *PostRepo) GetAll(ctx context.Context, limit, offset int) ([]*model.Post, error) {
	query := `
		SELECT id, title, content, author_id, created_at 
		FROM posts 
		ORDER BY created_at DESC 
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query all posts: %w", err)
	}
	defer rows.Close()

	posts := make([]*model.Post, 0, limit)
	for rows.Next() {
		var post model.Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan post row: %w", err)
		}
		posts = append(posts, &post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during posts rows iteration: %w", err)
	}

	return posts, nil
}

// Возвращает общее количество постов в системе
func (r *PostRepo) GetTotalCount(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM posts`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get total posts count: %w", err)
	}

	return count, nil
}

// Проверяет существование поста по его ID
func (r *PostRepo) Exists(ctx context.Context, id int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check post existence: %w", err)
	}

	return exists, nil
}

// Получает посты конкретного автора с поддержкой пагинации
func (r *PostRepo) GetByAuthorID(ctx context.Context, authorID int, limit, offset int) ([]*model.Post, error) {
	query := `
		SELECT id, title, content, author_id, created_at 
		FROM posts 
		WHERE author_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, authorID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query posts by author_id: %w", err)
	}
	defer rows.Close()

	posts := make([]*model.Post, 0, limit)
	for rows.Next() {
		var post model.Post
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan author post row: %w", err)
		}
		posts = append(posts, &post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error during author posts rows iteration: %w", err)
	}

	return posts, nil
}

// Возвращает общее количество постов конкретного автора
func (r *PostRepo) GetTotalCountByAuthorID(ctx context.Context, authorID int) (int, error) {
	query := `SELECT COUNT(*) FROM posts WHERE author_id = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, authorID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get author posts count: %w", err)
	}

	return count, nil
}
