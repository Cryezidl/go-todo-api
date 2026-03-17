package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
	Role         string    `db:"user_role"`
	Created_At   time.Time `db:"created_at"`
	Updated_At   time.Time `db:"updated_at"`
	Is_Active    bool      `db:"is_active"`
	Timezone     string    `db:"timezone"`
	Language     string    `db:"lang"`
	Theme        string    `db:"theme"`
}
