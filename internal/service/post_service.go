package service

import (
	"advanced-blog-management-system/internal/errors/apperrors"
	"advanced-blog-management-system/internal/model"
	"advanced-blog-management-system/internal/repository"
	"context"
	"errors"
	"fmt"
	"unicode/utf8"
)

type PostService struct {
	postRepo repository.PostRepository
	userRepo repository.UserRepository
}

// Создает новый экземпляр PostService.
// Сигнатура изменена: теперь возвращается интерфейс PostServiceInterface.
func NewPostService(postRepo repository.PostRepository, userRepo repository.UserRepository) PostServiceInterface {
	return &PostService{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

// Валидирует входящий запрос и создает новый пост
func (s *PostService) Create(ctx context.Context, userID int, req *model.PostCreateRequest) (*model.Post, error) {
	// 1. Дополнительная бизнес-валидация строк
	titleLen := utf8.RuneCountInString(req.Title)
	if titleLen == 0 || titleLen > 200 {
		return nil, errors.New("post title must be between 1 and 200 characters")
	}
	if utf8.RuneCountInString(req.Content) == 0 {
		return nil, errors.New("post content cannot be empty")
	}

	// 2. Сборка модели данных
	post := &model.Post{
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: userID,
	}

	// 3. Сохранение в репозиторий PostgreSQL
	if err := s.postRepo.Create(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to save post: %w", err)
	}

	return post, nil
}

// Возвращает пост по ID, если он существует
func (s *PostService) GetByID(ctx context.Context, id int, requestorID int) (*model.Post, error) {
	if id <= 0 {
		return nil, apperrors.ErrInvalidPostID
	}

	post, err := s.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query post by id: %w", err)
	}

	// Если репозиторий вернул nil — значит, запись в БД отсутствует
	if post == nil {
		return nil, apperrors.ErrPostNotFound
	}

	return post, nil
}

// Возвращает список всех постов с пагинацией и метаданными
func (s *PostService) GetAll(ctx context.Context, limit, offset int) ([]*model.Post, int, error) {
	// Валидация и корректировка параметров
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	// Получаем общее количество постов в системе
	total, err := s.postRepo.GetTotalCount(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total posts count: %w", err)
	}

	if total == 0 {
		return []*model.Post{}, 0, nil
	}

	// Выборка порции данных
	posts, err := s.postRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch posts: %w", err)
	}

	return posts, total, nil
}

// Возвращает публикации конкретного пользователя
func (s *PostService) GetByAuthor(ctx context.Context, authorID int, limit, offset int) ([]*model.Post, int, error) {
	if authorID <= 0 {
		return nil, 0, apperrors.ErrUserNotFound
	}

	// Проверяем параметры пагинации
	if limit <= 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	// Получаем количество постов конкретного автора
	total, err := s.postRepo.GetTotalCountByAuthorID(ctx, authorID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get author posts count: %w", err)
	}

	if total == 0 {
		return []*model.Post{}, 0, nil
	}

	// Запрашиваем посты автора из базы данных
	posts, err := s.postRepo.GetByAuthorID(ctx, authorID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch author posts: %w", err)
	}

	return posts, total, nil
}
