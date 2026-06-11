package service

import (
	"advanced-blog-management-system/internal/errors/apperrors"
	"advanced-blog-management-system/internal/model"
	"advanced-blog-management-system/internal/repository"
	"advanced-blog-management-system/pkg/auth"
	"context"
	"fmt"
)

type UserService struct {
	userRepo   repository.UserRepository
	jwtManager *auth.JWTManager // Используем внедренный менеджер для генерации токенов
}

// Создает новый экземпляр UserService.
// Сигнатура изменена: теперь возвращается интерфейс UserServiceInterface, а вместо строки принимает *auth.JWTManager
func NewUserService(userRepo repository.UserRepository, jwtManager *auth.JWTManager) UserServiceInterface {
	return &UserService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// Проверяет уникальность данных, валидирует сложность пароля, хэширует его и возвращает TokenResponse
func (s *UserService) Register(ctx context.Context, req *model.UserCreateRequest) (*model.TokenResponse, error) {
	// 1. Проверяем, не занят ли email
	emailExists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email availability: %w", err)
	}
	if emailExists {
		return nil, apperrors.ErrUserAlreadyExists
	}

	// 2. Проверяем, не занят ли username
	usernameExists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to check username availability: %w", err)
	}
	if usernameExists {
		return nil, apperrors.ErrUserAlreadyExists
	}

	// 3. Опциональная бизнес-проверка сложности пароля
	if err := auth.ValidatePasswordStrength(req.Password); err != nil {
		return nil, err
	}

	// 4. Безопасное хэширование пароля через нашу утилиту
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 5. Собираем модель пользователя для сохранения в PostgreSQL
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to persist new user: %w", err)
	}

	// 6. Генерируем токен с помощью обновленного метода объекта JWTManager
	token, expiresAt, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token after registration: %w", err)
	}

	return &model.TokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user.ToResponse(),
	}, nil
}

// Аутентифицирует пользователя по email и паролю, возвращая TokenResponse
func (s *UserService) Login(ctx context.Context, req *model.UserLoginRequest) (*model.TokenResponse, error) {
	// 1. Ищем пользователя по email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to query user during login: %w", err)
	}
	// Безопасное сокрытие деталей: при отсутствии юзера отдаем ErrInvalidCredentials
	if user == nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	// 2. Проверяем соответствие пароля хэшу через утилиту
	if !auth.CheckPassword(req.Password, user.Password) {
		return nil, apperrors.ErrInvalidCredentials
	}

	// 3. Генерируем токен через метод JWTManager
	token, expiresAt, err := s.jwtManager.GenerateToken(user.ID, user.Email, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token during login: %w", err)
	}

	return &model.TokenResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      user.ToResponse(),
	}, nil
}

// Находит пользователя по его ID
func (s *UserService) GetByID(ctx context.Context, id int) (*model.User, error) {
	if id <= 0 {
		return nil, apperrors.ErrUserNotFound
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user by id: %w", err)
	}
	if user == nil {
		return nil, apperrors.ErrUserNotFound
	}

	return user, nil
}
