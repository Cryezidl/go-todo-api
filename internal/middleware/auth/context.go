package auth

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey string

var userIDKey ctxKey = "UserID"

func GetUserId(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return id, ok
}
