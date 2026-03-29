package user

import (
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/Cryezidl/go-todo-api/pkg/myerrors"
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

var (
	allowedLanguages = map[string]bool{"RU": true, "ENG": true}
	allowedThemes    = map[string]bool{"light": true, "dark": true, "system": true}
	emailRegex       = regexp.MustCompile(`^[a-zA-Z0-9._%+]+@[a-zA-Z0-9]+\.[a-z]{2,}$`)
)

func (u *UpdateUserProfile) Validate() error {
	if u.Username != nil {
		if err := validateUsername(*u.Username); err != nil {
			return err
		}
	}

	if u.Language != nil {
		if !allowedLanguages[*u.Language] {
			return myerrors.ErrInvalidLanguage
		}
	}

	if u.Theme != nil {
		if !allowedThemes[*u.Theme] {
			return myerrors.ErrInvalidTheme
		}
	}

	if u.Timezone != nil {
		if _, err := time.LoadLocation(*u.Timezone); err != nil {
			return myerrors.ErrInvalidTimezone
		}
	}
	return nil
}

func (u *UpdateUserPassword) Validate() error {
	if len(u.NewPassword) < 8 || len(u.NewPassword) > 72 {
		return myerrors.ErrPasswordInvalidLength
	}
	if u.NewPassword == u.OldPassword {
		return myerrors.ErrPasswordSameAsOld
	}
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, r := range u.NewPassword {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	if b2i(hasUpper)+b2i(hasLower)+b2i(hasDigit)+b2i(hasSpecial) < 3 {
		return myerrors.ErrPasswordTooWeak
	}
	return nil
}

func (u *UpdateUserEmail) Validate() error {
	if len(u.Password) < 8 || len(u.Password) > 72 {
		return myerrors.ErrPasswordInvalidLength
	}
	u.NewEmail = strings.TrimSpace(u.NewEmail)
	if len(u.NewEmail) < 5 {
		return myerrors.ErrEmailTooShort
	}
	if !emailRegex.MatchString(u.NewEmail) {
		return myerrors.ErrEmailInvalid
	}
	return nil
}
func validateUsername(username string) error {
	if len(strings.TrimSpace(username)) <= 2 {
		return myerrors.ErrUsernameTooShort
	}
	containsLetters := false
	for _, s := range username {
		if unicode.IsLetter(s) {
			containsLetters = true
			break
		}
	}

	if !containsLetters {
		return myerrors.ErrUsernameNoLetters
	}
	return nil
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}
