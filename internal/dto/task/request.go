package task

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type TaskRequest struct {
	UserID      uuid.UUID  `json:"-"`
	ListID      uuid.UUID  `json:"list_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    int        `json:"priority"`
	Deadline    *time.Time `json:"deadline"`
	RemindAt    *time.Time `json:"remind_at"`
	Tags        []string   `json:"tags"`
}

type TaskUpdateRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Status      *string    `json:"status"`
	Deadline    *time.Time `json:"deadline"`
	Priority    *int       `json:"priority"`
	RemindAt    *time.Time `json:"remind_at"`
	Tags        *[]string  `json:"tags"`
}

type TaskUpdateStatusRequest struct {
	Status string `json:"status" validate:"required"`
}

var (
	availableStatus = map[string]bool{"done": true, "in_progress": true, "pending": true}
)

func (r *TaskUpdateRequest) Validate() error {
	// Title validation
	if r.Title != nil && *r.Title == "" {
		return errors.New("title cannot be empty")
	}

	// Status validation
	if r.Status != nil {
		if _, ok := availableStatus[*r.Status]; !ok {
			return errors.New("unavailable status")
		}
	}

	// RemindAt validation (with nil check!)
	if r.RemindAt != nil && r.RemindAt.Before(time.Now()) {
		return errors.New("unavailable date for reminder")
	}

	// Deadline validation (with nil check!)
	if r.Deadline != nil && r.Deadline.Before(time.Now()) {
		return errors.New("unavailable date for deadline")
	}

	// Priority validation
	if r.Priority != nil && *r.Priority < 0 {
		return errors.New("priority can't be a negative number")
	}

	return nil
}

func (r *TaskUpdateStatusRequest) Validate() error {
	if _, ok := availableStatus[r.Status]; !ok {
		return errors.New("unavailable status")
	}
	return nil
}

func (r *TaskRequest) Validate() error {
	// Title validation
	if r.Title == "" {
		return errors.New("title cannot be empty")
	}

	// Status validation
	if _, ok := availableStatus[r.Status]; !ok {
		return errors.New("unavailable status")
	}

	// RemindAt validation (with nil check!)
	if r.RemindAt != nil && r.RemindAt.Before(time.Now()) {
		return errors.New("unavailable date for reminder")
	}

	// Deadline validation (with nil check!)
	if r.Deadline != nil && r.Deadline.Before(time.Now()) {
		return errors.New("unavailable date for deadline")
	}

	// Priority validation
	if r.Priority < 0 {
		return errors.New("priority can't be a negative number")
	}

	return nil
}
