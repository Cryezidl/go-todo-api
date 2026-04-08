package tasklist

import "github.com/google/uuid"

type TaskListRequest struct {
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsPrivate   bool      `json:"is_private"`
}

type TaskListUpdateRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	IsPrivate   *bool   `json:"is_private"`
}
