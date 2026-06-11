package model

import (
	"github.com/go-playground/validator/v10"
	"time"
	"unicode"
)

// Глобальный экземпляр валидатора
var validate = validator.New()

func init() {
	// Регистрируем кастомный тег валидации "password_strength" для проверки сложности пароля
	_ = validate.RegisterValidation("password_strength", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		if len(password) < 8 {
			return false
		}

		var (
			hasUpper   bool
			hasLower   bool
			hasNumber  bool
			hasSpecial bool
		)

		for _, r := range password {
			switch {
			case unicode.IsUpper(r):
				hasUpper = true
			case unicode.IsLower(r):
				hasLower = true
			case unicode.IsNumber(r):
				hasNumber = true
			case unicode.IsPunct(r) || unicode.IsSymbol(r):
				hasSpecial = true
			}
		}

		return hasUpper && hasLower && hasNumber && hasSpecial
	})
}

// Представляет модель пользователя в системе
type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Представляет модель поста в блоге
type Post struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	AuthorID  int       `json:"author_id" db:"author_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Представляет модель комментария к посту
type Comment struct {
	ID        int       `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	PostID    int       `json:"post_id" db:"post_id"`
	AuthorID  int       `json:"author_id" db:"author_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Представляет запрос на создание пользователя
type UserCreateRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password_strength"`
}

// Представляет запрос на вход пользователя
type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Представляет запрос на создание поста
type PostCreateRequest struct {
	Title   string `json:"title" validate:"required,min=1,max=200"`
	Content string `json:"content" validate:"required,min=1"`
}

// Представляет запрос на создание комментария
type CommentCreateRequest struct {
	Content string `json:"content" validate:"required,min=1,max=1000"`
}

// Структура для ответа с данными пользователя
type UserResponse struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Структура для ответа с JWT токеном
type TokenResponse struct {
	Token     string       `json:"token"`
	ExpiresAt time.Time    `json:"expires_at"`
	User      UserResponse `json:"user"`
}

// Структура для ответа с данными поста
type PostResponse struct {
	ID        int          `json:"id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Author    UserResponse `json:"author"`
	CreatedAt time.Time    `json:"created_at"`
}

// Структура для ответа с данными комментария
type CommentResponse struct {
	ID        int          `json:"id"`
	Content   string       `json:"content"`
	PostID    int          `json:"post_id"`
	Author    UserResponse `json:"author"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// Преобразует User в UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

// Проверяет валидность полей запроса на регистрацию
func (r *UserCreateRequest) Validate() error {
	return validate.Struct(r)
}

// Проверяет валидность полей запроса на вход
func (r *UserLoginRequest) Validate() error {
	return validate.Struct(r)
}

// Проверяет валидность полей запроса на создание поста
func (r *PostCreateRequest) Validate() error {
	return validate.Struct(r)
}

// Проверяет валидность полей запроса на создание комментария
func (r *CommentCreateRequest) Validate() error {
	return validate.Struct(r)
}
