package auth

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	dtoAuth "github.com/Cryezidl/go-todo-api/internal/dto/auth"
	service "github.com/Cryezidl/go-todo-api/internal/service/auth"
	"github.com/Cryezidl/go-todo-api/pkg/httputils"
	"github.com/Cryezidl/go-todo-api/pkg/myerrors"
)

type AuthHandler struct {
	authService *service.AuthService
	log         *slog.Logger
	tokenDur    time.Duration
}

func NewAuthHandler(authService *service.AuthService, log *slog.Logger, tokenDur time.Duration) *AuthHandler {
	return &AuthHandler{authService: authService, log: log, tokenDur: tokenDur}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	//считать данные в RegisterInput
	input := dtoAuth.RegisterInput{}
	if err := httputils.ReadFromJSON(r, &input); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, myerrors.ErrInvalidBody.Error(), h.log)
		return
	}
	//валидировать данные
	if err := input.Validate(); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, err.Error(), h.log)
		return
	}
	//вызвать сервис
	resp, token, err := h.authService.RegisterUser(r.Context(), input)
	if err != nil {
		if errors.Is(err, myerrors.ErrEmailTaken) || errors.Is(err, myerrors.ErrUsernameTaken) {
			httputils.RespondWithError(w, http.StatusConflict, err.Error(), h.log)
			return
		}
		httputils.RespondWithError(w, http.StatusInternalServerError, myerrors.ErrInternalServer.Error(), h.log)
		return
	}
	//Засунуть jwt в хедер
	w.Header().Set("Authorization", "Bearer "+token)
	w.Header().Set("X-Expires-In", fmt.Sprintf("%.0f", h.tokenDur.Seconds()))
	//вернуть created
	httputils.RespondWithJSON(w, http.StatusCreated, resp, h.log)

}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	//беру LoginInput
	input := dtoAuth.LoginInput{}
	if err := httputils.ReadFromJSON(r, &input); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, myerrors.ErrInvalidBody.Error(), h.log)
		return
	}
	//валидирую
	if err := input.Validate(); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, err.Error(), h.log)
		return
	}
	//вызываю сервис
	resp, token, err := h.authService.Login(r.Context(), input)
	if err != nil {
		if errors.Is(err, myerrors.ErrInvalidCredentials) {
			httputils.RespondWithError(w, http.StatusUnauthorized, err.Error(), h.log)
			return
		}
		httputils.RespondWithError(w, http.StatusInternalServerError, myerrors.ErrInternalServer.Error(), h.log)
		return
	}
	//записываю токен
	w.Header().Set("Authorization", "Bearer "+token)
	w.Header().Set("X-Expires-In", fmt.Sprintf("%.0f", h.tokenDur.Seconds()))
	//вернуть created
	httputils.RespondWithJSON(w, http.StatusOK, resp, h.log)
}
