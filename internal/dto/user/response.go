package user

import (
	"time"

	"github.com/google/uuid"
)

type UserResponse struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	Role       string    `json:"role"`
	Created_At time.Time `json:"created_at"`
	Timezone   string    `json:"timezone"`
	Language   string    `json:"language"`
	Theme      string    `json:"theme"`
}

type PublicUserResponse struct {
	ID         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	Created_At time.Time `json:"created_at"`
}
