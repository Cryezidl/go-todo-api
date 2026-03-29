package validator

import (
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/Cryezidl/go-todo-api/pkg/myerrors"
)

var (
	allowedLanguages = map[string]bool{"RU": true, "ENG": true}
	allowedThemes    = map[string]bool{"light": true, "dark": true, "system": true}
	emailRegex       = regexp.MustCompile(`^[a-zA-Z0-9._%+]+@[a-zA-Z0-9]+\.[a-z]{2,}$`)
)

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if len(email) < 5 {
		return myerrors.ErrEmailTooShort
	}
	if !emailRegex.MatchString(email) {
		return myerrors.ErrEmailInvalid
	}
	return nil
}

func ValidateUsername(username string) error {
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

func ValidatePassword(password string) error {
	if len(password) < 8 || len(password) > 72 {
		return myerrors.ErrPasswordInvalidLength
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, r := range password {
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

func ValidateLanguage(lang string) error {
	if !allowedLanguages[strings.ToUpper(lang)] {
		return myerrors.ErrInvalidLanguage
	}
	return nil
}

func ValidateTheme(theme string) error {
	if !allowedThemes[strings.ToLower(theme)] {
		return myerrors.ErrInvalidTheme
	}
	return nil
}

func ValidateTimezone(tz string) error {
	if _, err := time.LoadLocation(tz); err != nil {
		return myerrors.ErrInvalidTimezone
	}
	return nil
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}
