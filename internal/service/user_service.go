package service

import (
	"advanced-blog-management-system/internal/errors/apperrors"
	"advanced-blog-management-system/internal/model"
	"advanced-blog-management-system/internal/repository"
	"advanced-blog-management-system/pkg/auth"
	"context"
	"errors"
	"fmt"
)

type UserService struct {
	userRepo   repository.UserRepository
	jwtManager *auth.JWTManager
}

func NewUserService(userRepo repository.UserRepository, jwtManager *auth.JWTManager) *UserService {
	return &UserService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// TODO: Реализовать Register()
// Валидировать, проверить существование email/username, захешировать пароль,
// создать пользователя, генерировать JWT токен, вернуть TokenResponse
func (s *UserService) Register(ctx context.Context, req *model.UserCreateRequest) (*model.TokenResponse, error) {
	// TODO: реализовать
	return nil, nil
}

// TODO: Реализовать Login()
// Валидировать, получить пользователя по email, проверить пароль,
// генерировать JWT токен, вернуть TokenResponse
func (s *UserService) Login(ctx context.Context, req *model.UserLoginRequest) (*model.TokenResponse, error) {
	// TODO: реализовать
	return nil, nil
}

// TODO: Реализовать GetByID()
// Получить пользователя из репозитория по ID
func (s *UserService) GetByID(ctx context.Context, id int) (*model.User, error) {
	// TODO: реализовать
	return nil, nil
}

// TODO: Реализовать GetByEmail()
// Получить пользователя из репозитория по email
func (s *UserService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	// TODO: реализовать
	return nil, nil
}
