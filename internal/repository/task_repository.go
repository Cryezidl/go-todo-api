package repository

import (
	"context"
	"time"

	"github.com/Cryezidl/go-todo-api/internal/model"
	"github.com/google/uuid"
)

type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) error
	FindById(ctx context.Context, ID uuid.UUID) (*model.Task, error)
	FindByListID(ctx context.Context, taskListID uuid.UUID) ([]*model.Task, error)
	Update(ctx context.Context, task *model.Task) error
	UpdateStatus(ctx context.Context, ID uuid.UUID, status string, completedAt *time.Time) error
	Delete(ctx context.Context, ID uuid.UUID) error
}
