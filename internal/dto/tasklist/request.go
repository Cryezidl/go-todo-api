package tasklist

import (
	"errors"

	"github.com/google/uuid"
)

type TaskListRequest struct {
	UserID      uuid.UUID `json:"-"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsPrivate   bool      `json:"is_private"`
}

type TaskListUpdateRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	IsPrivate   *bool   `json:"is_private"`
}

func (r *TaskListUpdateRequest) Validate() error {
	// Опциональная валидация
	if r.Title != nil && *r.Title == "" {
		return errors.New("title cannot be empty")
	}
	return nil
}
func (r *TaskListRequest) Validate() error {
	// Опциональная валидация
	if r.Title == "" {
		return errors.New("title cannot be empty")
	}
	return nil
}
