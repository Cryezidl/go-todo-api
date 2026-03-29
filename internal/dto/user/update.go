package user

import (
	"github.com/Cryezidl/go-todo-api/pkg/myerrors"
	"github.com/Cryezidl/go-todo-api/pkg/validator"
)

type UpdateUserProfile struct {
	Username *string `json:"username"`
	Timezone *string `json:"timezone"`
	Language *string `json:"language"`
	Theme    *string `json:"theme"`
}

type UpdateUserPassword struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=30"`
}

type UpdateUserEmail struct {
	Password string `json:"password" validate:"required,email,max=256"`
	NewEmail string `json:"new_email" validate:"required,min=8,max=30"`
}

func (u *UpdateUserProfile) Validate() error {
	if u.Username != nil {
		if err := validator.ValidateUsername(*u.Username); err != nil {
			return err
		}
	}

	if u.Language != nil {
		if err := validator.ValidateLanguage(*u.Language); err != nil {
			return err
		}
	}

	if u.Theme != nil {
		if err := validator.ValidateTheme(*u.Theme); err != nil {
			return err
		}
	}

	if u.Timezone != nil {
		if err := validator.ValidateTimezone(*u.Timezone); err != nil {
			return err
		}
	}
	return nil
}

func (u *UpdateUserPassword) Validate() error {
	if u.NewPassword == u.OldPassword {
		return myerrors.ErrPasswordSameAsOld
	}
	return validator.ValidatePassword(u.NewPassword)
}

func (u *UpdateUserEmail) Validate() error {
	if len(u.Password) < 8 {
		return myerrors.ErrPasswordInvalidLength
	}
	return validator.ValidateEmail(u.NewEmail)
}
