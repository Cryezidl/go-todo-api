package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateJWTToken(secretKey []byte, ID uuid.UUID, role string, dur time.Duration) (string, error) {
	claim := jwt.MapClaims{
		"sub":  ID.String(),
		"exp":  time.Now().Add(dur).Unix(),
		"role": role,
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(secretKey)
}
