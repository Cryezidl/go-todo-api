package task

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type TaskRequest struct {
	UserID      uuid.UUID `json:"-"`
	ListID      uuid.UUID `json:"list_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    int       `json:"priority"`
	RemindAt    time.Time `json:"remind_at"`
	Tags        []string  `json:"tags"`
}

type TaskUpdateRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Status      *string    `json:"status"`
	Priority    *int       `json:"priority"`
	RemindAt    *time.Time `json:"remind_at"`
	Tags        *[]string  `json:"tags"`
}

type TaskUpdateStatusRequest struct {
	Status string `json:"status" validate:"required"`
}

var (
	availableStatus = map[string]bool{"done": true, "in progress": true}
)

func (r *TaskUpdateRequest) Validate() error {
	// Опциональная валидация
	if r.Title != nil && *r.Title == "" {
		return errors.New("title cannot be empty")
	}
	if _, ok := availableStatus[*r.Status]; !ok {
		return errors.New("unavailable status")
	}
	if r.RemindAt.Before(time.Now()) {
		return errors.New("unavailable date for reminder")
	}

	if *r.Priority < 0 {
		return errors.New("priority can't be a negative number")
	}

	return nil
}

func (r *TaskUpdateStatusRequest) Validate() error {
	// Опциональная валидация
	if _, ok := availableStatus[r.Status]; !ok {
		return errors.New("unavailable status")
	}
	return nil
}

func (r *TaskRequest) Validate() error {
	// Опциональная валидация
	if r.Title == "" {
		return errors.New("title cannot be empty")
	}
	if _, ok := availableStatus[r.Status]; !ok {
		return errors.New("unavailable status")
	}
	if r.RemindAt.Before(time.Now()) {
		return errors.New("unavailable date for reminder")
	}
	if r.Priority < 0 {
		return errors.New("priority can't be a negative number")
	}
	return nil
}
