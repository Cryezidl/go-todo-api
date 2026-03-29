package user

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	dtoUser "github.com/Cryezidl/go-todo-api/internal/dto/user"
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

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	//взять id из контекста
	id, ok := auth.GetUserId(r.Context())
	if !ok {
		httputils.RespondWithError(w, http.StatusUnauthorized, myerrors.ErrUnauthorized.Error(), h.log)
		return
	}
	//считать данные в UpdateUserProfile
	updateData := dtoUser.UpdateUserProfile{}
	if err := httputils.ReadFromJSON(r, &updateData); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, myerrors.ErrInvalidBody.Error(), h.log)
		return
	}
	//валидировать данные
	if err := updateData.Validate(); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, err.Error(), h.log)
		return
	}
	//вызвать метод из сервиса
	userResp, err := h.userService.UpdateProfile(r.Context(), updateData, id)
	if err != nil {
		if errors.Is(err, myerrors.ErrUserNotFound) {
			httputils.RespondWithError(w, http.StatusNotFound, err.Error(), h.log)
			return
		}
		httputils.RespondWithError(w, http.StatusInternalServerError, myerrors.ErrInternalServer.Error(), h.log)
		return
	}
	//вернуть обновленные данные UserResponse
	httputils.RespondWithJSON(w, http.StatusOK, userResp, h.log)
}

func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	//взять Id из контекста
	id, ok := auth.GetUserId(r.Context())
	if !ok {
		httputils.RespondWithError(w, http.StatusUnauthorized, myerrors.ErrUnauthorized.Error(), h.log)
		return
	}
	//считать данные с запроса
	updateData := dtoUser.UpdateUserPassword{}
	if err := httputils.ReadFromJSON(r, &updateData); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, myerrors.ErrInvalidBody.Error(), h.log)
		return
	}
	//провалидировать пароль
	if err := updateData.Validate(); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, err.Error(), h.log)
		return
	}
	//вызвать сервис
	if err := h.userService.UpdatePassword(r.Context(), updateData, id); err != nil {
		if errors.Is(err, myerrors.ErrUserNotFound) {
			httputils.RespondWithError(w, http.StatusNotFound, err.Error(), h.log)
			return
		}
		if errors.Is(err, myerrors.ErrInvalidCredentials) {
			httputils.RespondWithError(w, http.StatusForbidden, myerrors.ErrInvalidPassword.Error(), h.log)
			return
		}
		httputils.RespondWithError(w, http.StatusInternalServerError, myerrors.ErrInternalServer.Error(), h.log)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) UpdateEmail(w http.ResponseWriter, r *http.Request) {
	//взять Id из контекста
	id, ok := auth.GetUserId(r.Context())
	if !ok {
		httputils.RespondWithError(w, http.StatusUnauthorized, myerrors.ErrUnauthorized.Error(), h.log)
		return
	}
	//считать данные с запроса
	updateData := dtoUser.UpdateUserEmail{}
	if err := httputils.ReadFromJSON(r, &updateData); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, myerrors.ErrInvalidBody.Error(), h.log)
		return
	}

	//провалидировать почту
	if err := updateData.Validate(); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, err.Error(), h.log)
		return
	}

	//вызвать сервис
	if err := h.userService.UpdateEmail(r.Context(), updateData, id); err != nil {
		if errors.Is(err, myerrors.ErrUserNotFound) {
			httputils.RespondWithError(w, http.StatusNotFound, err.Error(), h.log)
			return
		}
		if errors.Is(err, myerrors.ErrInvalidCredentials) {
			httputils.RespondWithError(w, http.StatusForbidden, myerrors.ErrInvalidPassword.Error(), h.log)
			return
		}
		if errors.Is(err, myerrors.ErrEmailTaken) {
			httputils.RespondWithError(w, http.StatusConflict, err.Error(), h.log)
			return
		}
		httputils.RespondWithError(w, http.StatusInternalServerError, myerrors.ErrInternalServer.Error(), h.log)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request)
