package service

import (
	"advanced-blog-management-system/internal/model"
	"advanced-blog-management-system/internal/repository"
	"context"
	"fmt"
)

type PostService struct {
	postRepo repository.PostRepository
	userRepo repository.UserRepository
}

func NewPostService(postRepo repository.PostRepository, userRepo repository.UserRepository) *PostService {
	return &PostService{
		postRepo: postRepo,
		userRepo: userRepo,
	}
}

// TODO: Реализовать Create()
// Валидировать, создать пост, сохранить в БД, вернуть Post
func (s *PostService) Create(ctx context.Context, userID int, req *model.PostCreateRequest) (*model.Post, error) {
	// TODO: реализовать
	return nil, nil
}

// TODO: Реализовать GetByID()
// Получить пост по ID, вернуть Post или ошибку
func (s *PostService) GetByID(ctx context.Context, id int, requestorID int) (*model.Post, error) {
	// TODO: реализовать
	return nil, nil
}

// TODO: Реализовать GetAll()
// Валидировать limit/offset, получить посты, вернуть слайс и общее количество
func (s *PostService) GetAll(ctx context.Context, limit, offset int) ([]*model.Post, int, error) {
	// TODO: реализовать
	return nil, 0, nil
}

// TODO: Реализовать GetByAuthor()
// Валидировать limit/offset, получить посты автора, вернуть слайс и общее количество
func (s *PostService) GetByAuthor(ctx context.Context, authorID int, limit, offset int) ([]*model.Post, int, error) {
	// TODO: реализовать
	return nil, 0, nil
}
