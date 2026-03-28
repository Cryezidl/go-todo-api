package repository

import (
	"context"

	"github.com/Cryezidl/go-todo-api/internal/model"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByUserId(ctx context.Context, id uuid.UUID) (*model.User, error)
	FindByUserEmail(ctx context.Context, email string) (*model.User, error)
	FindByUserName(ctx context.Context, name string) (*model.User, error)
	FindByUserEmailOrName(ctx context.Context, login string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
