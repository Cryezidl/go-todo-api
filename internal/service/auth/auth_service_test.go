package auth

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	dtoAuth "github.com/Cryezidl/go-todo-api/internal/dto/auth"
	"github.com/Cryezidl/go-todo-api/internal/model"
	"github.com/Cryezidl/go-todo-api/internal/service/user/mock"
)

var testLogger *slog.Logger

func TestAuthService_RegisterUser(t *testing.T) {
	tests := []struct {
		name      string
		mockError error
		req       dtoAuth.RegisterInput
		checkUser func(*testing.T, model.User)
		wantErr   bool
	}{
		{
			name:      "create with all fields",
			mockError: nil,
			req: dtoAuth.RegisterInput{
				Email:    "alice@example.com",
				Username: "Alice",
				Password: "correctpassword",
				Timezone: "Moscow",
				Language: "RU",
				Theme:    "Dark",
			},
			wantErr: false,
			checkUser: func(t *testing.T, us model.User) {
				if us.Email != "alice@example.com" {
					t.Errorf("Email = %v, want alice@example.com", us.Email)
				}
				if us.Username != "Alice" {
					t.Errorf("Username = %v, want Alice", us.Username)
				}
				if us.Timezone != "Moscow" {
					t.Errorf("Timezone = %v, want Moscow", us.Timezone)
				}
				if us.Language != "RU" {
					t.Errorf("Language = %v, want RU", us.Language)
				}
				if us.Theme != "Dark" {
					t.Errorf("Theme = %v, want Dark", us.Theme)
				}

			},
		},
		{
			name:      "create with only required fields",
			mockError: nil,
			req: dtoAuth.RegisterInput{
				Email:    "alice@example.com",
				Username: "Alice",
				Password: "correctpassword",
			},
			wantErr: false,
			checkUser: func(t *testing.T, us model.User) {
				if us.Email != "alice@example.com" {
					t.Errorf("Email = %v, want alice@example.com", us.Email)
				}
				if us.Username != "Alice" {
					t.Errorf("Username = %v, want Alice", us.Username)
				}
				if us.Timezone != "UTC" {
					t.Errorf("Timezone = %v, want UTC", us.Timezone)
				}
				if us.Language != "ENG" {
					t.Errorf("Language = %v, want ENG", us.Language)
				}
				if us.Theme != "system" {
					t.Errorf("Theme = %v, want system", us.Theme)
				}

			},
		},
		{
			name:      "repository error",
			mockError: errors.New("database error"),
			req: dtoAuth.RegisterInput{
				Password: "correctpassword",
			},
			checkUser: nil,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock.MockUserRepository{}
			mockRepo.CreateError = tt.mockError

			s := NewAuthService(mockRepo, testLogger, "secretkey", time.Minute)

			_, token, err := s.RegisterUser(context.Background(), tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			mockRepo.AssertCreate(t, 1)

			if !tt.wantErr {

				if tt.checkUser != nil && mockRepo.LastCreateUser != nil {
					tt.checkUser(t, *mockRepo.LastCreateUser)
				}

				if token == "" {
					t.Error("token should not be empty")
				}
			}
		})

	}
}
