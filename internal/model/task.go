package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Task struct {
	ID          uuid.UUID      `db:"id"`
	ListID      uuid.UUID      `db:"list_id"`
	UserID      uuid.UUID      `db:"user_id"`
	Title       string         `db:"title"`
	Description string         `db:"description"`
	Status      string         `db:"status"`
	Priority    int            `db:"priority"`
	Deadline    *time.Time     `db:"deadline"`
	RemindAt    *time.Time     `db:"remind_at"`
	Tags        pq.StringArray `db:"tags"`
	CompletedAt *time.Time     `db:"completed_at"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}
