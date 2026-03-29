package myerrors

import (
	"errors"
)

var (
	//servie or db
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailTaken         = errors.New("email already registered")
	ErrUsernameTaken      = errors.New("username already registered")
	ErrInvalidCredentials = errors.New("invalid credentials")

	//handlers
	ErrUnauthorized   = errors.New("unauthorized")
	ErrInternalServer = errors.New("internal server error")
	ErrInvalidBody    = errors.New("invalid request body")

	//middleware
	ErrMissingAuth  = errors.New("missing authorization header")
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")

	//validation
	ErrUsernameTooShort      = errors.New("username must be at least 3 characters long")
	ErrUsernameNoLetters     = errors.New("username must contain at least one letter")
	ErrEmailInvalid          = errors.New("invalid email format")
	ErrEmailTooShort         = errors.New("email must be at least 5 characters long")
	ErrInvalidTheme          = errors.New("Invalid theme: supported values are light, dark, system")
	ErrInvalidLanguage       = errors.New("invalid language: supported values are ru, en")
	ErrInvalidTimezone       = errors.New("invalid timezone: must be a valid IANA location (e.g. Europe/Moscow)")
	ErrPasswordInvalidLength = errors.New("password must be between 8 and 72 characters")
	ErrPasswordSameAsOld     = errors.New("new password must be different from the old one")
	ErrPasswordTooWeak       = errors.New("password must include upper, lower case letters and digits")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrInvalidLoginFormat    = errors.New("invalid login format")
)
