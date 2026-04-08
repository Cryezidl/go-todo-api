package model

import (
	"time"

	"github.com/google/uuid"
)

type TaskList struct {
	ID          uuid.UUID `db:"id"`
	UserID      uuid.UUID `db:"user_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	IsPrivate   bool      `db:"is_private"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// func (u *User) ToResponse() dtoUser.UserResponse {
// 	return dtoUser.UserResponse{
// 		ID:         u.ID,
// 		Email:      u.Email,
// 		Username:   u.Username,
// 		Role:       u.Role,
// 		Created_At: u.Created_At,
// 		Timezone:   u.Timezone,
// 		Language:   u.Language,
// 		Theme:      u.Theme,
// 	}
// }

// func (u *User) ToPublicResponse() dtoUser.PublicUserResponse {
// 	return dtoUser.PublicUserResponse{
// 		ID:         u.ID,
// 		Username:   u.Username,
// 		Created_At: u.Created_At,
// 	}
// }
