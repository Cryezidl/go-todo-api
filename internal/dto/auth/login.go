package auth

import (
	"github.com/Cryezidl/go-todo-api/pkg/myerrors"
	"github.com/Cryezidl/go-todo-api/pkg/validator"
)

type LoginInput struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=30"`
}

func (l *LoginInput) Validate() error {
	errUsername := validator.ValidateUsername(l.Login)
	errEmail := validator.ValidateEmail(l.Login)

	if errUsername != nil && errEmail != nil {
		return myerrors.ErrInvalidLoginFormat
	}
	if len(l.Password) < 8 || len(l.Password) > 72 {
		return myerrors.ErrPasswordInvalidLength
	}
	return nil
}
