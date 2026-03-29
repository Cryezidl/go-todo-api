package jwtutils

import (
	"errors"
	"fmt"
	"time"

	"github.com/Cryezidl/go-todo-api/pkg/myerrors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomClaims struct {
	Role                 string `json:"role"`
	jwt.RegisteredClaims        // Здесь уже есть Subject (ID), ExpiresAt, IssuedAt
}

func GenerateJWTToken(secretKey []byte, ID uuid.UUID, role string, dur time.Duration) (string, error) {
	claims := CustomClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   ID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(dur)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ParseJWTToken(tokenStr string, secretKey []byte) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		// Проверяем, не истек ли срок действия
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, myerrors.ErrExpiredToken
		}
		return nil, myerrors.ErrInvalidToken
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, myerrors.ErrInvalidToken
}
