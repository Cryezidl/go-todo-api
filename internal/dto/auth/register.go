package auth

import "github.com/Cryezidl/go-todo-api/pkg/validator"

type RegisterInput struct {
	Email    string `json:"email" validate:"required,email,max=256"`
	Username string `json:"username" validate:"required,min=2,max=30"`
	Password string `json:"password" validate:"required,min=8,max=30"`
	Timezone string `json:"timezone"`
	Language string `json:"language"`
	Theme    string `json:"theme"`
}

func (r *RegisterInput) Validate() error {
	if err := validator.ValidateEmail(r.Email); err != nil {
		return err
	}

	if err := validator.ValidateUsername(r.Username); err != nil {
		return err
	}

	if err := validator.ValidatePassword(r.Password); err != nil {
		return err
	}

	if err := validator.ValidateTimezone(r.Timezone); err != nil {
		return err
	}

	if err := validator.ValidateLanguage(r.Language); err != nil {
		return err
	}

	if err := validator.ValidateTheme(r.Theme); err != nil {
		return err
	}
	return nil
}
