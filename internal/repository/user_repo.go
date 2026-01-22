package repository

import (
	"advanced-blog-management-system/internal/model"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

// TODO: Реализовать Create()
// INSERT INTO users (username, email, password, created_at, updated_at) VALUES (...) RETURNING id
// Установить CreatedAt, UpdatedAt = time.Now(), сканировать ID в user
func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	// TODO: реализовать
	return nil
}

// TODO: Реализовать GetByID()
// SELECT ... FROM users WHERE id = $1
// Если no rows - вернуть nil, nil; если error - вернуть fmt.Errorf()
func (r *UserRepo) GetByID(ctx context.Context, id int) (*model.User, error) {
	// TODO: реализовать
	return nil, nil
}

// TODO: Реализовать GetByEmail()
// SELECT ... FROM users WHERE email = $1
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	// TODO: реализовать
	return nil, nil
}

// TODO: Реализовать GetByUsername()
// SELECT ... FROM users WHERE username = $1
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	// TODO: реализовать
	return nil, nil
}

// TODO: Реализовать ExistsByEmail()
// SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
func (r *UserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	// TODO: реализовать
	return false, nil
}

// TODO: Реализовать ExistsByUsername()
// SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)
func (r *UserRepo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	// TODO: реализовать
	return false, nil
}

// TODO: Реализовать Update()
// UPDATE users SET username = $1, email = $2, password = $3, updated_at = $4 WHERE id = $5
// Установить UpdatedAt = time.Now(), проверить RowsAffected, вернуть ошибку если 0 rows
func (r *UserRepo) Update(ctx context.Context, user *model.User) error {
	// TODO: реализовать
	return nil
}

// TODO: Реализовать Delete()
// DELETE FROM users WHERE id = $1
// Проверить RowsAffected, вернуть ошибку если 0 rows
func (r *UserRepo) Delete(ctx context.Context, id int) error {
	// TODO: реализовать
	return nil
}
