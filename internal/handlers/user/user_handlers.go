package user

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/Cryezidl/go-todo-api/internal/middleware/auth"
	service "github.com/Cryezidl/go-todo-api/internal/service/user"
	"github.com/Cryezidl/go-todo-api/pkg/httputils"
	"github.com/Cryezidl/go-todo-api/pkg/myerrors"
)

type UserHandler struct {
	userService *service.UserService
	log         *slog.Logger
}

func NewUserHandler(userService *service.UserService, log *slog.Logger) *UserHandler {
	return &UserHandler{userService: userService, log: log}
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	//извлечь id
	id, ok := auth.GetUserId(r.Context())
	if !ok {
		httputils.RespondWithError(w, http.StatusUnauthorized, myerrors.ErrUnauthorized.Error(), h.log)
		return
	}

	//найти по айдишнику профиль
	user, err := h.userService.GetById(r.Context(), id)
	if err != nil {
		if errors.Is(err, myerrors.ErrUserNotFound) {
			httputils.RespondWithError(w, http.StatusNotFound, err.Error(), h.log)
			return
		}
		httputils.RespondWithError(w, http.StatusInternalServerError, myerrors.ErrInternalServer.Error(), h.log)
		return
	}
	//вернуть профиль
	httputils.RespondWithJSON(w, http.StatusOK, user, h.log)
}

func (h *UserHandler) GetByName(w http.ResponseWriter, r *http.Request) {
	//извлекаем из строки запроса ник
	username := chi.URLParam(r, "username")
	if strings.TrimSpace(username) == "" {
		httputils.RespondWithError(w, http.StatusNotFound, myerrors.ErrUserNotFound.Error(), h.log)
		return
	}
	//получаем юзера
	user, err := h.userService.GetByName(r.Context(), username)
	if err != nil {
		if errors.Is(err, myerrors.ErrUserNotFound) {
			httputils.RespondWithError(w, http.StatusNotFound, err.Error(), h.log)
			return
		}
		httputils.RespondWithError(w, http.StatusInternalServerError, myerrors.ErrInternalServer.Error(), h.log)
		return
	}
	//возвращаем
	httputils.RespondWithJSON(w, http.StatusOK, user, h.log)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request)

func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request)

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request)
