package user

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	dtoUser "github.com/Cryezidl/go-todo-api/internal/dto/user"
	"github.com/Cryezidl/go-todo-api/internal/model"
	"github.com/Cryezidl/go-todo-api/internal/service/user/mock"
	"github.com/Cryezidl/go-todo-api/pkg/hash"
	"github.com/Cryezidl/go-todo-api/pkg/myerrors"
	"github.com/google/uuid"
)

func ptr[T any](v T) *T {
	return &v
}

var testLogger *slog.Logger

func init() {
	// Тихий логгер для тестов (ничего не выводит)
	testLogger = slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
}

func TestUserService_GetById(t *testing.T) {
	userID := uuid.New()
	expectedUser := &model.User{
		ID:       userID,
		Username: "Alice",
	}

	tests := []struct {
		name       string
		id         uuid.UUID
		mockResult *model.User
		mockError  error
		wantId     uuid.UUID
		wantName   string
		wantErr    bool
	}{
		{
			name:       "success",
			id:         userID,
			mockResult: expectedUser,
			mockError:  nil,
			wantErr:    false,
			wantId:     userID,
			wantName:   "Alice",
		},
		{
			name:       "user not found",
			id:         uuid.New(),
			mockResult: nil,
			mockError:  myerrors.ErrUserNotFound,
			wantErr:    true,
			wantId:     uuid.Nil,
			wantName:   "",
		},
		{
			name:       "repository error",
			id:         userID,
			mockResult: nil,
			mockError:  errors.New("database error"),
			wantErr:    true,
			wantId:     uuid.Nil,
			wantName:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &mock.MockUserRepository{}
			mockRepo.FindByUserIdResult = tt.mockResult
			mockRepo.FindByUserIdError = tt.mockError

			s := NewUserService(mockRepo, testLogger)

			got, gotErr := s.GetById(context.Background(), tt.id)

			if (gotErr != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", gotErr, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.ID != userID {
					t.Errorf("ID = %v, want %v", got.ID, userID)
				}
				if got.Username != "Alice" {
					t.Errorf("Username = %s, want Alice", got.Username)
				}
			}

			mockRepo.AssertFindByUserId(t, 1)
		})
	}
}

func TestUserService_GetByEmail(t *testing.T) {
	userID := uuid.New()
	email := "alice@example.com"

	expectedUser := &model.User{
		ID:    userID,
		Email: email,
	}

	tests := []struct {
		name       string
		email      string
		mockResult *model.User
		mockError  error
		wantErr    bool
		wantEmail  string
		wantId     uuid.UUID
	}{
		{
			name:       "success",
			email:      email,
			mockResult: expectedUser,
			mockError:  nil,
			wantErr:    false,
			wantEmail:  email,
			wantId:     userID,
		},
		{
			name:       "user not fount",
			email:      "wrong@example.com",
			mockResult: nil,
			mockError:  myerrors.ErrUserNotFound,
			wantErr:    true,
			wantEmail:  "",
			wantId:     uuid.Nil,
		},
		{
			name:       "repository error",
			email:      email,
			mockResult: nil,
			mockError:  errors.New("database error"),
			wantErr:    true,
			wantEmail:  "",
			wantId:     uuid.Nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &mock.MockUserRepository{}
			mockRepo.FindByUserEmailResult = tt.mockResult
			mockRepo.FindByUserEmailError = tt.mockError

			s := NewUserService(mockRepo, testLogger)

			got, gotErr := s.GetByEmail(context.Background(), tt.email)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", gotErr, tt.wantErr)
				return
			}
			// TODO: update the condition below to compare got with tt.want.
			if !tt.wantErr {
				if tt.wantEmail != got.Email {
					t.Errorf("email = %v, want %v", got.Email, tt.wantEmail)
				}
				if tt.wantId != got.ID {
					t.Errorf("ID = %v, want %v", got.ID, tt.wantId)
				}
			}

			mockRepo.AssertFindByUserEmail(t, 1)
		})
	}
}

func TestUserService_GetByName(t *testing.T) {
	userID := uuid.New()
	username := "alice"

	expectedUser := &model.User{
		ID:       userID,
		Username: username,
	}

	tests := []struct {
		name         string
		username     string
		mockResult   *model.User
		mockError    error
		wantErr      bool
		wantUsername string
		wantId       uuid.UUID
	}{
		{
			name:         "success",
			username:     username,
			mockResult:   expectedUser,
			mockError:    nil,
			wantErr:      false,
			wantUsername: username,
			wantId:       userID,
		},
		{
			name:         "user not fount",
			username:     "john",
			mockResult:   nil,
			mockError:    myerrors.ErrUserNotFound,
			wantErr:      true,
			wantUsername: "",
			wantId:       uuid.Nil,
		},
		{
			name:         "repository error",
			username:     username,
			mockResult:   nil,
			mockError:    errors.New("database error"),
			wantErr:      true,
			wantUsername: "",
			wantId:       uuid.Nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &mock.MockUserRepository{}
			mockRepo.FindByUserNameResult = tt.mockResult
			mockRepo.FindByUserNameError = tt.mockError

			s := NewUserService(mockRepo, testLogger)

			got, gotErr := s.GetByName(context.Background(), tt.username)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", gotErr, tt.wantErr)
				return
			}
			// TODO: update the condition below to compare got with tt.want.
			if !tt.wantErr {
				if tt.wantUsername != got.Username {
					t.Errorf("name = %v, want %v", got.Username, tt.wantUsername)
				}
				if tt.wantId != got.ID {
					t.Errorf("ID = %v, want %v", got.ID, tt.wantId)
				}
			}

			mockRepo.AssertFindByUserName(t, 1)
		})
	}
}

func TestUserService_UpdateProfile(t *testing.T) {
	userID := uuid.New()
	existingUser := &model.User{
		ID:       userID,
		Username: "oldname",
		Email:    "old@example.com",
		Language: "ENG",
		Timezone: "UTC",
		Theme:    "light",
	}

	tests := []struct {
		name            string
		id              uuid.UUID
		newData         dtoUser.UpdateUserProfile
		mockFindResult  *model.User
		mockFindError   error
		mockUpdateError error
		wantErr         bool
		wantUpdateCall  int
		checkUser       func(*testing.T, *model.User)
	}{
		{
			name: "update all fields",
			id:   userID,
			newData: dtoUser.UpdateUserProfile{
				Username: ptr("newname"),
				Language: ptr("RU"),
				Timezone: ptr("Europe/Moscow"),
				Theme:    ptr("dark"),
			},
			mockFindResult:  existingUser,
			mockFindError:   nil,
			mockUpdateError: nil,
			wantErr:         false,
			checkUser: func(t *testing.T, u *model.User) {
				if u.Username != "newname" {
					t.Errorf("Username = %s, want newname", u.Username)
				}
				if u.Language != "RU" {
					t.Errorf("Language = %s, want RU", u.Language)
				}
				if u.Timezone != "Europe/Moscow" {
					t.Errorf("Timezone = %s, want Europe/Moscow", u.Timezone)
				}
				if u.Theme != "dark" {
					t.Errorf("Theme = %s, want dark", u.Theme)
				}
				if u.Email != existingUser.Email {
					t.Errorf("Email changed to %s, but should stay old@example.com", u.Email)
				}
				if u.ID != existingUser.ID {
					t.Errorf("ID changed to %s, but should stay %s", u.ID.String(), existingUser.ID.String())
				}
			},
			wantUpdateCall: 1,
		},
		{
			name: "update only username",
			id:   userID,
			newData: dtoUser.UpdateUserProfile{
				Username: ptr("newname"),
			},
			mockFindResult:  existingUser,
			mockFindError:   nil,
			mockUpdateError: nil,
			wantErr:         false,
			checkUser: func(t *testing.T, u *model.User) {
				if u.Username != "newname" {
					t.Errorf("Username = %s, want newname", u.Username)
				}
				if u.Language != existingUser.Language {
					t.Errorf("Language = %s, want %s", u.Language, existingUser.Language)
				}
				if u.Timezone != existingUser.Timezone {
					t.Errorf("Timezone = %s, want %s", u.Timezone, existingUser.Timezone)
				}
				if u.Theme != existingUser.Theme {
					t.Errorf("Theme = %s, want %s", u.Theme, existingUser.Theme)
				}
				if u.Email != existingUser.Email {
					t.Errorf("Email changed to %s, but should stay old@example.com", u.Email)
				}
				if u.ID != existingUser.ID {
					t.Errorf("ID changed to %s, but should stay %s", u.ID.String(), existingUser.ID.String())
				}
			},
			wantUpdateCall: 1,
		},
		{
			name: "user not found",
			id:   uuid.New(),
			newData: dtoUser.UpdateUserProfile{
				Username: ptr("newname"),
			},
			mockFindResult:  nil,
			mockFindError:   myerrors.ErrUserNotFound,
			mockUpdateError: nil,
			wantErr:         true,
			checkUser:       nil,
			wantUpdateCall:  0,
		},
		{
			name: "update repository error",
			id:   userID,
			newData: dtoUser.UpdateUserProfile{
				Username: ptr("newname"),
			},
			mockFindResult:  existingUser,
			mockUpdateError: errors.New("database error"),
			wantErr:         true,
			checkUser:       nil,
			wantUpdateCall:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &mock.MockUserRepository{}
			mockRepo.FindByUserIdResult = tt.mockFindResult
			mockRepo.UpdateError = tt.mockUpdateError

			s := NewUserService(mockRepo, testLogger)

			got, gotErr := s.UpdateProfile(context.Background(), tt.newData, tt.id)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", gotErr, tt.wantErr)
				return
			}

			mockRepo.AssertFindByUserId(t, 1)
			mockRepo.AssertUpdate(t, tt.wantUpdateCall)
			if !tt.wantErr {
				if tt.checkUser != nil && mockRepo.LastUpdateUser != nil {
					tt.checkUser(t, mockRepo.LastUpdateUser)
				}
				if got.ID != tt.id {
					t.Errorf("response ID = %v, want %v", got.ID, tt.id)
				}
			}

		})
	}
}

/*
	type UpdateUserEmail struct {
		Password string `json:"password" validate:"required,email,max=256"`
		NewEmail string `json:"new_email" validate:"required,min=8,max=30"`
	}
*/
func TestUserService_UpdateEmail(t *testing.T) {
	userID := uuid.New()
	hashPassword := func(password string) string {
		hash, _ := hash.HashPassword(password)
		return hash
	}

	existingUser := &model.User{
		ID:           userID,
		PasswordHash: hashPassword("correctpassword"),
		Email:        "old@example.com",
	}

	tests := []struct {
		name                string
		id                  uuid.UUID
		newData             dtoUser.UpdateUserEmail
		mockFindResult      *model.User
		mockFindError       error
		mockFindEmailResult *model.User
		mockFindEmailError  error
		mockUpdateError     error
		wantErr             bool
		wantFindEmailCall   int
		checkUser           func(*testing.T, *model.User)
	}{
		{
			name: "success",
			id:   userID,
			newData: dtoUser.UpdateUserEmail{
				Password: "correctpassword",
				NewEmail: "new@example.com",
			},
			mockFindResult:    existingUser,
			wantErr:           false,
			mockFindError:     nil,
			wantFindEmailCall: 1,
			checkUser: func(t *testing.T, u *model.User) {
				if u.Email != "new@example.com" {
					t.Errorf("Email = %s, want new@example.com", u.Email)
				}
				if u.PasswordHash != existingUser.PasswordHash {
					t.Errorf("PasswordHash changed to %s, but should stay %s", u.PasswordHash, existingUser.PasswordHash)
				}
			},
		},
		{
			name: "user not found",
			id:   uuid.New(),
			newData: dtoUser.UpdateUserEmail{
				Password: "correctpassword",
				NewEmail: "new@example.com",
			},
			mockFindError:     myerrors.ErrUserNotFound,
			wantErr:           true,
			wantFindEmailCall: 0,
		},
		{
			name: "wrong password",
			id:   userID,
			newData: dtoUser.UpdateUserEmail{
				Password: "wrongpassword",
				NewEmail: "new@example.com",
			},
			mockFindResult:    existingUser,
			wantErr:           true,
			wantFindEmailCall: 0,
		},
		{
			name: "email already taken",
			id:   userID,
			newData: dtoUser.UpdateUserEmail{
				Password: "correctpassword",
				NewEmail: "new@example.com",
			},
			mockFindResult:     existingUser,
			wantErr:            true,
			wantFindEmailCall:  1,
			mockFindEmailError: myerrors.ErrEmailTaken,
		},
		{
			name: "update repository error",
			id:   userID,
			newData: dtoUser.UpdateUserEmail{
				Password: "correctpassword",
				NewEmail: "new@example.com",
			},
			mockFindResult:      existingUser,
			wantErr:             true,
			mockFindEmailResult: &model.User{Email: "new@example.com"},
			wantFindEmailCall:   1,
			mockUpdateError:     errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &mock.MockUserRepository{}
			mockRepo.FindByUserIdResult = tt.mockFindResult
			mockRepo.FindByUserIdError = tt.mockFindError
			mockRepo.FindByUserEmailError = tt.mockFindEmailError
			mockRepo.FindByUserEmailResult = tt.mockFindEmailResult
			mockRepo.UpdateError = tt.mockUpdateError

			s := NewUserService(mockRepo, testLogger)

			gotErr := s.UpdateEmail(context.Background(), tt.newData, tt.id)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", gotErr, tt.wantErr)
				return
			}

			mockRepo.AssertFindByUserId(t, 1)
			mockRepo.AssertFindByUserEmail(t, tt.wantFindEmailCall)

			if tt.wantErr {
				mockRepo.AssertUpdate(t, 0)
			} else {
				mockRepo.AssertUpdate(t, 1)
				if tt.checkUser != nil && mockRepo.LastUpdateUser != nil {
					tt.checkUser(t, mockRepo.LastUpdateUser)
				}
			}
		})
	}
}

/*
	type UpdateUserPassword struct {
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=8,max=30"`
	}
*/
func TestUserService_UpdatePassword(t *testing.T) {
	userID := uuid.New()
	hashPassword := func(password string) string {
		hash, _ := hash.HashPassword(password)
		return hash
	}
	existingUser := &model.User{
		ID:           userID,
		PasswordHash: hashPassword("oldpassword"),
	}
	oldHash := existingUser.PasswordHash

	tests := []struct {
		name            string
		id              uuid.UUID
		newData         dtoUser.UpdateUserPassword
		mockFindResult  *model.User
		mockFindError   error
		mockUpdateError error
		wantErr         bool
		checkUser       func(*testing.T, *model.User)
	}{
		{
			name: "success",
			id:   userID,
			newData: dtoUser.UpdateUserPassword{
				OldPassword: "oldpassword",
				NewPassword: "newpassword",
			},
			mockFindResult: existingUser,
			wantErr:        false,
			checkUser: func(t *testing.T, u *model.User) {
				if u.PasswordHash == oldHash {
					t.Errorf("PasswordHash %s, should be changed", u.PasswordHash)
				}
				if u.ID != existingUser.ID {
					t.Errorf("ID changed to %s, but should stay %s", u.ID.String(), existingUser.ID.String())
				}
			},
		},
		{
			name: "user not found",
			id:   uuid.New(),
			newData: dtoUser.UpdateUserPassword{
				OldPassword: "oldpassword",
				NewPassword: "newpassword",
			},
			mockFindResult: nil,
			mockFindError:  myerrors.ErrUserNotFound,
			wantErr:        true,
		},
		{
			name: "wrong password",
			id:   userID,
			newData: dtoUser.UpdateUserPassword{
				OldPassword: "wrongpassword",
				NewPassword: "newpassword",
			},
			mockFindResult: existingUser,
			wantErr:        true,
		},
		{
			name: "repository update error",
			id:   userID,
			newData: dtoUser.UpdateUserPassword{
				OldPassword: "oldpassword",
				NewPassword: "newpassword",
			},
			mockFindResult:  existingUser,
			mockUpdateError: errors.New("database error"),
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &mock.MockUserRepository{}
			mockRepo.FindByUserIdResult = tt.mockFindResult
			mockRepo.FindByUserIdError = tt.mockFindError
			mockRepo.UpdateError = tt.mockUpdateError

			s := NewUserService(mockRepo, testLogger)

			gotErr := s.UpdatePassword(context.Background(), tt.newData, tt.id)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", gotErr, tt.wantErr)
				return
			}

			mockRepo.AssertFindByUserId(t, 1)
			if tt.wantErr {
				mockRepo.AssertUpdate(t, 0)
			} else {
				mockRepo.AssertUpdate(t, 1)
				if tt.checkUser != nil && mockRepo.LastUpdateUser != nil {
					tt.checkUser(t, mockRepo.LastUpdateUser)
				}
			}
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	userID := uuid.New()
	existingUser := &model.User{
		ID:       userID,
		Username: "Alice",
	}

	tests := []struct {
		name            string
		id              uuid.UUID
		mockFindResult  *model.User
		mockFindError   error
		mockDeleteError error
		wantErr         bool
	}{
		{
			name:           "success",
			id:             userID,
			mockFindResult: existingUser,
			wantErr:        false,
		},
		{
			name:           "user not found",
			id:             uuid.New(),
			mockFindResult: nil,
			mockFindError:  myerrors.ErrUserNotFound,
			wantErr:        true,
		},
		{
			name:            "repository delete error",
			id:              userID,
			mockFindResult:  existingUser,
			mockDeleteError: errors.New("database error"),
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &mock.MockUserRepository{}
			mockRepo.FindByUserIdResult = tt.mockFindResult
			mockRepo.FindByUserIdError = tt.mockFindError
			mockRepo.DeleteError = tt.mockDeleteError

			s := NewUserService(mockRepo, testLogger)

			gotErr := s.Delete(context.Background(), tt.id)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", gotErr, tt.wantErr)
				return
			}
			mockRepo.AssertFindByUserId(t, 1)
			if tt.wantErr == false || tt.mockDeleteError != nil {
				mockRepo.AssertDelete(t, 1)
			} else {
				mockRepo.AssertDelete(t, 0)
			}

		})
	}
}
