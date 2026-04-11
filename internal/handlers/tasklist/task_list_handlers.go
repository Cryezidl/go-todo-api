package tasklist

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	dtoTaskList "github.com/Cryezidl/go-todo-api/internal/dto/tasklist"
	"github.com/Cryezidl/go-todo-api/internal/middleware/auth"
	"github.com/Cryezidl/go-todo-api/internal/model"
	service "github.com/Cryezidl/go-todo-api/internal/service/tasklist"
	"github.com/Cryezidl/go-todo-api/pkg/httputils"
	"github.com/Cryezidl/go-todo-api/pkg/myerrors"
)

type TaskListHandlers struct {
	taskListService service.TaskListService
	log             *slog.Logger
}

func NewTaskListHandlers(taskListService service.TaskListService, log *slog.Logger) *TaskListHandlers {
	return &TaskListHandlers{taskListService: taskListService, log: log}
}

func (h *TaskListHandlers) GetAll(w http.ResponseWriter, r *http.Request) {
	//извлечь id
	userId, ok := auth.GetUserId(r.Context())
	if !ok {
		httputils.RespondWithError(w, http.StatusUnauthorized, myerrors.ErrUnauthorized.Error(), h.log)
		return
	}

	taskLists, err := h.taskListService.GetByUserID(r.Context(), userId)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, myerrors.ErrInternalServer.Error(), h.log)
		return
	}
	if taskLists == nil {
		taskLists = []*model.TaskList{}
	}

	httputils.RespondWithJSON(w, http.StatusOK, taskLists, h.log)
}

func (h *TaskListHandlers) GetById(w http.ResponseWriter, r *http.Request) {
	//извлечь user id
	userID, ok := auth.GetUserId(r.Context())
	if !ok {
		httputils.RespondWithError(w, http.StatusUnauthorized, myerrors.ErrUnauthorized.Error(), h.log)
		return
	}

	//излвечь аргумент
	idParam := chi.URLParam(r, "id")
	taskListID, err := uuid.Parse(idParam)
	if err != nil {
		httputils.RespondWithError(w, http.StatusNotFound, myerrors.ErrTaskListNotFound.Error(), h.log)
		return
	}

	taskList, err := h.taskListService.GetByID(r.Context(), taskListID)
	if err != nil {
		if errors.Is(err, myerrors.ErrTaskListNotFound) {
			httputils.RespondWithError(w, http.StatusNotFound, err.Error(), h.log)
			return
		}
		httputils.RespondWithError(w, http.StatusInternalServerError, myerrors.ErrInternalServer.Error(), h.log)
		return
	}

	if taskList.UserID != userID && taskList.IsPrivate {
		httputils.RespondWithError(w, http.StatusForbidden, myerrors.ErrAccessDenied.Error(), h.log)
		return
	}

	httputils.RespondWithJSON(w, http.StatusOK, taskList, h.log)
}

func (h *TaskListHandlers) Update(w http.ResponseWriter, r *http.Request) {
	//извлечь user id
	userID, ok := auth.GetUserId(r.Context())
	if !ok {
		httputils.RespondWithError(w, http.StatusUnauthorized, myerrors.ErrUnauthorized.Error(), h.log)
		return
	}

	//собираем новые данные
	updateData := dtoTaskList.TaskListUpdateRequest{}
	if err := httputils.ReadFromJSON(r, &updateData); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, myerrors.ErrInvalidBody.Error(), h.log)
		return
	}

	//валидировать данные
	if err := updateData.Validate(); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, err.Error(), h.log)
		return
	}

	//излвечь аргумент
	idParam := chi.URLParam(r, "id")
	taskListID, err := uuid.Parse(idParam)
	if err != nil {
		httputils.RespondWithError(w, http.StatusNotFound, myerrors.ErrTaskListNotFound.Error(), h.log)
		return
	}

	//вызвать метод из сервиса
	updatedTasKList, err := h.taskListService.Update(r.Context(), updateData, taskListID, userID)
	if err != nil {
		if errors.Is(err, myerrors.ErrTaskListNotFound) {
			httputils.RespondWithError(w, http.StatusNotFound, err.Error(), h.log)
			return
		}
		if errors.Is(err, myerrors.ErrAccessDenied) {
			httputils.RespondWithError(w, http.StatusForbidden, err.Error(), h.log)
			return
		}
		httputils.RespondWithError(w, http.StatusInternalServerError, myerrors.ErrInternalServer.Error(), h.log)
		return
	}

	//вернуть обновленные данные UserResponse
	httputils.RespondWithJSON(w, http.StatusOK, updatedTasKList, h.log)
}

func (h *TaskListHandlers) Create(w http.ResponseWriter, r *http.Request) {
	//извлечь user id
	userID, ok := auth.GetUserId(r.Context())
	if !ok {
		httputils.RespondWithError(w, http.StatusUnauthorized, myerrors.ErrUnauthorized.Error(), h.log)
		return
	}

	//собираем новые данные
	createData := dtoTaskList.TaskListRequest{}
	if err := httputils.ReadFromJSON(r, &createData); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, myerrors.ErrInvalidBody.Error(), h.log)
		return
	}
	createData.UserID = userID

	//валидировать данные
	if err := createData.Validate(); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, err.Error(), h.log)
		return
	}

	//вызвать метод из сервиса
	createdTaskList, err := h.taskListService.Create(r.Context(), createData)
	if err != nil {
		if errors.Is(err, myerrors.ErrTaskListAlreadyExists) {
			httputils.RespondWithError(w, http.StatusConflict, err.Error(), h.log)
			return
		}
		httputils.RespondWithError(w, http.StatusInternalServerError, myerrors.ErrInternalServer.Error(), h.log)
		return
	}

	//вернуть обновленные данные
	w.Header().Set("Location", fmt.Sprintf("/api/task-lists/%s", createdTaskList.ID.String()))
	httputils.RespondWithJSON(w, http.StatusCreated, createdTaskList, h.log)
}

func (h *TaskListHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	//извлечь user id
	userID, ok := auth.GetUserId(r.Context())
	if !ok {
		httputils.RespondWithError(w, http.StatusUnauthorized, myerrors.ErrUnauthorized.Error(), h.log)
		return
	}

	//излвечь аргумент
	idParam := chi.URLParam(r, "id")
	taskListID, err := uuid.Parse(idParam)
	if err != nil {
		httputils.RespondWithError(w, http.StatusNotFound, myerrors.ErrTaskListNotFound.Error(), h.log)
		return
	}

	//вызвать метод из сервиса
	if err := h.taskListService.Delete(r.Context(), taskListID, userID); err != nil {
		if errors.Is(err, myerrors.ErrTaskListNotFound) {
			httputils.RespondWithError(w, http.StatusNotFound, err.Error(), h.log)
			return
		}
		if errors.Is(err, myerrors.ErrAccessDenied) {
			httputils.RespondWithError(w, http.StatusForbidden, err.Error(), h.log)
			return
		}
		httputils.RespondWithError(w, http.StatusInternalServerError, myerrors.ErrInternalServer.Error(), h.log)
		return
	}

	//вернуть обновленные данные
	w.WriteHeader(http.StatusNoContent)
}
