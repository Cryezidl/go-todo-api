package model

import (
	"time"

	dtoUser "github.com/Cryezidl/go-todo-api/internal/dto/user"
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

func (u *User) ToResponse() dtoUser.UserResponse {
	return dtoUser.UserResponse{
		ID:         u.ID,
		Email:      u.Email,
		Username:   u.Username,
		Role:       u.Role,
		Created_At: u.Created_At,
		Timezone:   u.Timezone,
		Language:   u.Language,
		Theme:      u.Theme,
	}
}
