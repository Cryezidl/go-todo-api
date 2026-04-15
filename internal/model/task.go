package model

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID  `db:"id"`
	ListID      uuid.UUID  `db:"list_id"`
	UserID      uuid.UUID  `db:"user_id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Status      string     `db:"status"`
	Priority    int        `db:"priority"`
	RemindAt    time.Time  `db:"remind_at"`
	Tags        []string   `db:"tags"`
	CompletedAt *time.Time `db:"completed_at"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}
