package user

import (
	"log/slog"

	service "github.com/Cryezidl/go-todo-api/internal/service/user"
)

type UserHadnler struct {
	userService *service.UserService
	log         *slog.Logger
}

func NewUserHander(userService *service.UserService, log *slog.Logger) *UserHadnler {
	return &UserHadnler{userService: userService, log: log}
}
