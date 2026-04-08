package user

import (
	"context"
	"log/slog"

	"errors"

	dtoUser "github.com/Cryezidl/go-todo-api/internal/dto/user"
	"github.com/Cryezidl/go-todo-api/internal/repository"
	"github.com/Cryezidl/go-todo-api/pkg/hash"
	"github.com/Cryezidl/go-todo-api/pkg/myerrors"
	"github.com/google/uuid"
)

type UserService struct {
	userRepository repository.UserRepository
	log            *slog.Logger
}

func NewUserService(userRepository repository.UserRepository, log *slog.Logger) *UserService {
	return &UserService{userRepository: userRepository, log: log}
}

func (s *UserService) GetById(ctx context.Context, id uuid.UUID) (dtoUser.UserResponse, error) {
	//const op = "service.user.UserService.GetById"

	if err := ctx.Err(); err != nil {
		return dtoUser.UserResponse{}, err
	}

	user, err := s.userRepository.FindByUserId(ctx, id)
	if err != nil {
		return dtoUser.UserResponse{}, err
	}
	if user == nil {
		s.log.Debug("user not found in service", slog.String("id", id.String()))
		return dtoUser.UserResponse{}, myerrors.ErrUserNotFound
	}
	return user.ToResponse(), nil
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (dtoUser.UserResponse, error) {
	if err := ctx.Err(); err != nil {
		return dtoUser.UserResponse{}, err
	}

	user, err := s.userRepository.FindByUserEmail(ctx, email)
	if err != nil {
		return dtoUser.UserResponse{}, err
	}
	if user == nil {
		s.log.Debug("user not found in service", slog.String("email", email))
		return dtoUser.UserResponse{}, myerrors.ErrUserNotFound
	}
	return user.ToResponse(), nil
}

func (s *UserService) GetByName(ctx context.Context, username string) (dtoUser.PublicUserResponse, error) {
	if err := ctx.Err(); err != nil {
		return dtoUser.PublicUserResponse{}, err
	}

	user, err := s.userRepository.FindByUserName(ctx, username)
	if err != nil {
		return dtoUser.PublicUserResponse{}, err
	}
	if user == nil {
		s.log.Debug("user not found in service", slog.String("username", username))
		return dtoUser.PublicUserResponse{}, myerrors.ErrUserNotFound
	}
	return user.ToPublicResponse(), nil
}

func (s *UserService) UpdateProfile(ctx context.Context, newData dtoUser.UpdateUserProfile, id uuid.UUID) (dtoUser.UserResponse, error) {
	const op = "service.UserService.UpdateProfile"

	if err := ctx.Err(); err != nil {
		return dtoUser.UserResponse{}, err
	}
	//Получить текущее состояние юзера model.User
	user, err := s.userRepository.FindByUserId(ctx, id)
	if err != nil {
		return dtoUser.UserResponse{}, err
	}
	if user == nil {
		s.log.Warn("user not found for update", slog.String("op", op), slog.String("id", id.String()))
		return dtoUser.UserResponse{}, myerrors.ErrUserNotFound
	}
	//заменяем данные этого объекта юзера
	if newData.Username != nil {
		user.Username = *newData.Username
	}
	if newData.Language != nil {
		user.Language = *newData.Language
	}
	if newData.Timezone != nil {
		user.Timezone = *newData.Timezone
	}
	if newData.Theme != nil {
		user.Theme = *newData.Theme
	}
	//вызываем метод репозитория
	err = s.userRepository.Update(ctx, user)
	if err != nil {
		return dtoUser.UserResponse{}, err
	}
	s.log.Info("profile updated successfully", slog.String("id", id.String()))
	return user.ToResponse(), nil
}

func (s *UserService) UpdateEmail(ctx context.Context, newData dtoUser.UpdateUserEmail, id uuid.UUID) error {
	const op = "service.UserService.UpdateEmail"

	if err := ctx.Err(); err != nil {
		return err
	}
	//Получить текущее состояние юзера model.User
	user, err := s.userRepository.FindByUserId(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		s.log.Warn("user not found for update", slog.String("op", op), slog.String("id", id.String()))
		return myerrors.ErrUserNotFound
	}

	//Проверка пароля
	if err := hash.CheckPassword(user.PasswordHash, newData.Password); err != nil {
		s.log.Debug("wrong password attempt", slog.String("op", op), slog.String("id", id.String()))
		return myerrors.ErrInvalidCredentials
	}

	// Проверка уникальности новой почты
	alreadyExists, err := s.userRepository.FindByUserEmail(ctx, newData.NewEmail)
	if err != nil && !errors.Is(err, myerrors.ErrUserNotFound) {
		return err
	}

	if alreadyExists != nil {
		return myerrors.ErrEmailTaken
	}

	//заменяем данные этого объекта юзера
	user.Email = newData.NewEmail

	//вызываем метод репозитория
	err = s.userRepository.Update(ctx, user)
	if err != nil {
		return err
	}
	s.log.Info("email updated successfully", slog.String("id", id.String()))
	return nil
}

func (s *UserService) UpdatePassword(ctx context.Context, newData dtoUser.UpdateUserPassword, id uuid.UUID) error {
	const op = "service.UserService.UpdatePassword"

	if err := ctx.Err(); err != nil {
		return err
	}
	//Получить текущее состояние юзера model.User
	user, err := s.userRepository.FindByUserId(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		s.log.Warn("user not found for update", slog.String("op", op), slog.String("id", id.String()))
		return myerrors.ErrUserNotFound
	}

	//Проверка старого пароля
	if err := hash.CheckPassword(user.PasswordHash, newData.OldPassword); err != nil {
		s.log.Debug("wrong password attempt", slog.String("op", op), slog.String("id", id.String()))
		return myerrors.ErrInvalidCredentials
	}

	//хешируем пароль
	newHash, err := hash.HashPassword(newData.NewPassword)
	if err != nil {
		s.log.Error("failed to hash password", slog.String("op", op), slog.Any("error", err))
		return err
	}

	//заменяем данные этого объекта юзера
	user.PasswordHash = newHash

	//вызываем метод репозитория
	err = s.userRepository.Update(ctx, user)
	if err != nil {
		return err
	}
	s.log.Info("password updated successfully", slog.String("id", id.String()))
	return nil
}

func (s *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "service.UserService.Delete"

	if err := ctx.Err(); err != nil {
		return err
	}

	//Получить текущее состояние юзера model.User
	user, err := s.userRepository.FindByUserId(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		s.log.Warn("user not found for delete", slog.String("op", op), slog.String("id", id.String()))
		return myerrors.ErrUserNotFound
	}

	//вызываем метод удаления
	if err := s.userRepository.Delete(ctx, id); err != nil {
		if errors.Is(err, myerrors.ErrUserNotFound) {
			return myerrors.ErrUserNotFound
		}
		return err
	}
	s.log.Info("user was deleted", slog.String("op", op), slog.String("id", id.String()))
	return nil
}
