package auth

import (
	"context"
	"log/slog"
	"time"

	dtoAuth "github.com/Cryezidl/go-todo-api/internal/dto/auth"
	dtoUser "github.com/Cryezidl/go-todo-api/internal/dto/user"
	"github.com/Cryezidl/go-todo-api/internal/repository"

	"github.com/Cryezidl/go-todo-api/internal/model"
	"github.com/Cryezidl/go-todo-api/pkg/hash"
	"github.com/Cryezidl/go-todo-api/pkg/jwt"
	"github.com/Cryezidl/go-todo-api/pkg/myerrors"
)

type AuthService struct {
	userRepository repository.UserRepository
	log            *slog.Logger
	Secret         string
	Expiration     time.Duration
}

func NewAuthService(userRepository repository.UserRepository, log *slog.Logger) *AuthService {
	return &AuthService{userRepository: userRepository, log: log}
}

func (s *AuthService) RegisterUser(ctx context.Context, req dtoAuth.RegisterInput) (dtoUser.UserResponse, error) {

	//хеширование пароля
	passwordHash, err := hash.HashPassword(req.Password)
	if err != nil {
		return dtoUser.UserResponse{}, err
	}

	//заполнение незаполненных полей
	if req.Theme == "" {
		req.Theme = "system"
	}
	if req.Language == "" {
		req.Language = "ENG"
	}
	if req.Timezone == "" {
		req.Language = "UTC"
	}
	if req.Role == "" {
		req.Role = "client"
	}

	//создаем user
	user := &model.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: passwordHash,
		Role:         req.Role,
		Timezone:     req.Timezone,
		Language:     req.Language,
		Theme:        req.Theme,
	}
	//create репозитория
	if err := s.userRepository.Create(ctx, user); err != nil {
		return dtoUser.UserResponse{}, err
	}
	//создать первый туду лист

	return user.ToResponse(), nil
}

func (s *AuthService) Login(ctx context.Context, req dtoAuth.LoginInput) (dtoUser.UserResponse, string, error) {
	//чекаем email и username
	user, err := s.userRepository.FindByUserEmailOrName(ctx, req.Login)
	if err != nil {
		return dtoUser.UserResponse{}, "", myerrors.ErrInvalidCredentials
	}

	//сравнить хеш
	if err := hash.CheckPassword(user.PasswordHash, req.Password); err != nil {
		return dtoUser.UserResponse{}, "", myerrors.ErrInvalidCredentials
	}
	//если ок то делаем jwt
	jwtKey, err := jwt.GenerateJWTToken([]byte(s.Secret), user.ID, user.Role, s.Expiration)
	if err != nil {
		return dtoUser.UserResponse{}, "", err
	}
	//отправляем jwt и инфу о юзере
	return user.ToResponse(), jwtKey, nil
}
