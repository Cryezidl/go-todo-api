package auth

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/Cryezidl/go-todo-api/pkg/httputils"
	"github.com/Cryezidl/go-todo-api/pkg/jwtutils"
	"github.com/Cryezidl/go-todo-api/pkg/myerrors"
)

func AuthMiddleware(secretKey string, log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//Получить заголовок
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				httputils.RespondWithError(w, http.StatusUnauthorized, myerrors.ErrMissingAuth.Error(), log)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				httputils.RespondWithError(w, http.StatusUnauthorized, myerrors.ErrInvalidToken.Error(), log)
				return
			}
			token := parts[0]

			claims, err := jwtutils.ParseJWTToken(token, []byte(secretKey))
			if err != nil {
				if errors.Is(err, myerrors.ErrExpiredToken) {
					httputils.RespondWithError(w, http.StatusUnauthorized, myerrors.ErrExpiredToken.Error(), log)
					return
				}
				httputils.RespondWithError(w, http.StatusUnauthorized, myerrors.ErrInvalidToken.Error(), log)
				return
			}

			sub, err := claims.GetSubject()
			if err != nil {
				httputils.RespondWithError(w, http.StatusUnauthorized, "invalid subject", log)
				return
			}
			userID, err := uuid.Parse(sub)
			if err != nil {
				httputils.RespondWithError(w, http.StatusUnauthorized, "invalid user id format", log)
				return
			}

			ctx := WithUserID(r.Context(), userID)
			ctx = WithUserRole(ctx, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
