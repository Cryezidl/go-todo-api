package repository

import (
	"context"

	"github.com/Cryezidl/go-todo-api/internal/model"
	"github.com/google/uuid"
)

type TaskListRepository interface {
	Create(ctx context.Context, taskList *model.TaskList) error
	FindById(ctx context.Context, ID uuid.UUID) (*model.TaskList, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.TaskList, error)
	Update(ctx context.Context, taskList *model.TaskList) error
	Delete(ctx context.Context, ID uuid.UUID) error
}
