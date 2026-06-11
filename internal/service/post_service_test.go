package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"advanced-blog-management-system/internal/errors/apperrors"
	"advanced-blog-management-system/internal/model"
)

// Реализует интерфейс repository.PostRepository для тестов
type mockPostRepository struct {
	posts      map[int]*model.Post
	shouldFail bool
}

func (m *mockPostRepository) Create(ctx context.Context, post *model.Post) error {
	if m.shouldFail {
		return errors.New("database error")
	}
	post.ID = len(m.posts) + 1
	post.CreatedAt = time.Now()
	m.posts[post.ID] = post
	return nil
}

func (m *mockPostRepository) GetByID(ctx context.Context, id int) (*model.Post, error) {
	if m.shouldFail {
		return nil, errors.New("database error")
	}
	post, exists := m.posts[id]
	if !exists {
		return nil, nil // По контракту чистой архитектуры: не найдено — возвращаем nil, nil
	}
	return post, nil
}

func (m *mockPostRepository) GetAll(ctx context.Context, limit, offset int) ([]*model.Post, error) {
	return nil, nil
}

func (m *mockPostRepository) GetTotalCount(ctx context.Context) (int, error) {
	return len(m.posts), nil
}

func (m *mockPostRepository) Exists(ctx context.Context, id int) (bool, error) {
	_, exists := m.posts[id]
	return exists, nil
}

func (m *mockPostRepository) GetByAuthorID(ctx context.Context, authorID int, limit, offset int) ([]*model.Post, error) {
	return nil, nil
}

func (m *mockPostRepository) GetTotalCountByAuthorID(ctx context.Context, authorID int) (int, error) {
	return 0, nil
}

// Пустая заглушка для выполнения контракта конструктора
type mockUserRepository struct{}

func (m *mockUserRepository) Create(ctx context.Context, user *model.User) error { return nil }
func (m *mockUserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	return nil, nil
}
func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, nil
}
func (m *mockUserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	return nil, nil
}
func (m *mockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return false, nil
}
func (m *mockUserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	return false, nil
}
func (m *mockUserRepository) Update(ctx context.Context, user *model.User) error { return nil }
func (m *mockUserRepository) Delete(ctx context.Context, id int) error           { return nil }

// Тестирование успешного создания поста
func TestPostService_Create_Success(t *testing.T) {
	mockRepo := &mockPostRepository{posts: make(map[int]*model.Post)}
	mockUserRepo := &mockUserRepository{}
	service := NewPostService(mockRepo, mockUserRepo)

	req := &model.PostCreateRequest{
		Title:   "Valid Title",
		Content: "This is a valid post content description text.",
	}

	post, err := service.Create(context.Background(), 1, req)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if post.ID != 1 {
		t.Errorf("expected post ID to be 1, got %d", post.ID)
	}
	if post.Title != req.Title {
		t.Errorf("expected title %s, got %s", req.Title, post.Title)
	}
}

// Тестирование валидации короткого заголовка (Бизнес-правило)
func TestPostService_Create_ValidationError(t *testing.T) {
	mockRepo := &mockPostRepository{posts: make(map[int]*model.Post)}
	mockUserRepo := &mockUserRepository{}
	service := NewPostService(mockRepo, mockUserRepo)

	req := &model.PostCreateRequest{
		Title:   "", // Пустой заголовок вызовет ошибку
		Content: "Valid Content",
	}

	_, err := service.Create(context.Background(), 1, req)

	if err == nil {
		t.Fatal("expected validation error for empty title, got nil")
	}
}

// Тестирование получения поста по ID, когда он отсутствует в базе (Обработка ошибок)
func TestPostService_GetByID_NotFound(t *testing.T) {
	mockRepo := &mockPostRepository{posts: make(map[int]*model.Post)}
	mockUserRepo := &mockUserRepository{}
	service := NewPostService(mockRepo, mockUserRepo)

	_, err := service.GetByID(context.Background(), 999, 1) // ID 999 нет в базе

	if !errors.Is(err, apperrors.ErrPostNotFound) {
		t.Errorf("expected ErrPostNotFound, got %v", err)
	}
}
