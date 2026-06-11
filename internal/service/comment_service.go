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

type CommentService struct {
	repo     repository.CommentRepository
	postRepo repository.PostRepository
}

// Создает новый экземпляр CommentService, принимая интерфейсы
func NewCommentService(repo repository.CommentRepository, postRepo repository.PostRepository) CommentServiceInterface {
	return &CommentService{
		repo:     repo,
		postRepo: postRepo,
	}
}

// Проверяет существование поста, валидирует данные и создает новый комментарий
func (s *CommentService) Create(ctx context.Context, userID, postID int, content string) (*model.Comment, error) {
	// 1. Валидация входного идентификатора поста
	if postID <= 0 {
		return nil, apperrors.ErrInvalidPostID
	}

	// 2. Бизнес-проверка: существует ли пост, который пользователь хочет прокомментировать
	postExists, err := s.postRepo.Exists(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to check post existence: %w", err)
	}
	if !postExists {
		return nil, apperrors.ErrPostNotFound
	}

	// 3. Валидация контента комментария (учитываем utf-8 символы)
	contentLen := utf8.RuneCountInString(content)
	if contentLen == 0 {
		return nil, errors.New("comment content cannot be empty")
	}
	if contentLen > 1000 {
		return nil, errors.New("comment content cannot exceed 1000 characters")
	}

	// 4. Формирование модели для слоя данных
	comment := &model.Comment{
		Content:  content,
		PostID:   postID,
		AuthorID: userID,
	}

	// 5. Сохранение в базу данных
	if err := s.repo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("failed to save comment to repository: %w", err)
	}

	return comment, nil
}

// Проверяет существование поста и возвращает список комментариев с пагинацией и общим количеством
func (s *CommentService) GetByPost(ctx context.Context, postID, limit, offset int) ([]*model.Comment, int, error) {
	// 1. Валидация идентификатора поста
	if postID <= 0 {
		return nil, 0, apperrors.ErrInvalidPostID
	}

	// 2. Проверяем, существует ли пост
	postExists, err := s.postRepo.Exists(ctx, postID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to check post existence: %w", err)
	}
	if !postExists {
		return nil, 0, apperrors.ErrPostNotFound
	}

	// 3. Валидация и корректировка параметров пагинации (бизнес-правила)
	if limit <= 0 {
		limit = 10 // Значение по умолчанию
	} else if limit > 100 {
		limit = 100 // Ограничение сверху
	}

	if offset < 0 {
		offset = 0
	}

	// 4. Получаем общее количество комментариев для этого поста
	total, err := s.repo.GetCountByPostID(ctx, postID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get comments count: %w", err)
	}

	// Если комментариев вообще нет, возвращаем пустой слайс вместо nil
	if total == 0 {
		return []*model.Comment{}, 0, nil
	}

	// 5. Выборка порции данных из репозитория
	comments, err := s.repo.GetByPostID(ctx, postID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to retrieve comments slice: %w", err)
	}

	return comments, total, nil
}
