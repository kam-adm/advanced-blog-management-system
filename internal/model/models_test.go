package model

import "testing"

// Проверяет работу go-playground/validator для регистрации с учетом сложности пароля
func TestUserCreateRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     UserCreateRequest
		wantErr bool
	}{
		{
			name:    "Valid request data with strong password",
			req:     UserCreateRequest{Username: "user", Email: "user@mail.ru", Password: "SecurePassword123!"},
			wantErr: false,
		},
		{
			name:    "Missing username",
			req:     UserCreateRequest{Username: "", Email: "user@mail.ru", Password: "SecurePassword123!"},
			wantErr: true,
		},
		{
			name:    "Invalid email format",
			req:     UserCreateRequest{Username: "user", Email: "not-an-email", Password: "SecurePassword123!"},
			wantErr: true,
		},
		{
			name:    "Password too short",
			req:     UserCreateRequest{Username: "user", Email: "user@mail.ru", Password: "Q1w!"},
			wantErr: true,
		},
		{
			name:    "Password missing special characters",
			req:     UserCreateRequest{Username: "user", Email: "user@mail.ru", Password: "password123"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("UserCreateRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
