package task

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	dtoTask "github.com/Cryezidl/go-todo-api/internal/dto/task"
	"github.com/Cryezidl/go-todo-api/internal/middleware/auth"
	"github.com/Cryezidl/go-todo-api/internal/model"
	service "github.com/Cryezidl/go-todo-api/internal/service/task"
	taskListService "github.com/Cryezidl/go-todo-api/internal/service/tasklist"
	"github.com/Cryezidl/go-todo-api/pkg/httputils"
	"github.com/Cryezidl/go-todo-api/pkg/myerrors"
)

type TaskHandlers struct {
	taskService     *service.TaskService
	taskListService *taskListService.TaskListService
	log             *slog.Logger
}

func NewTaskHandlers(taskService *service.TaskService, taskListService *taskListService.TaskListService, log *slog.Logger) *TaskHandlers {
	return &TaskHandlers{taskService: taskService, taskListService: taskListService, log: log}
}

// /tasks?list_id={id}
func (h *TaskHandlers) GetAll(w http.ResponseWriter, r *http.Request) {
	//извлечь id
	userID, ok := auth.GetUserId(r.Context())
	if !ok {
		httputils.RespondWithError(w, http.StatusUnauthorized, myerrors.ErrUnauthorized.Error(), h.log)
		return
	}

	//излвечь аргумент
	listIDParam := r.URL.Query().Get("list_id")
	if listIDParam == "" {
		httputils.RespondWithError(w, http.StatusBadRequest, "list_id is required", h.log)
		return
	}

	taskListID, err := uuid.Parse(listIDParam)
	if err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, "invalid list_id", h.log)
		return
	}

	//найти список
	tasklist, err := h.taskListService.GetByID(r.Context(), taskListID)
	if err != nil {
		if errors.Is(err, myerrors.ErrTaskListNotFound) {
			httputils.RespondWithError(w, http.StatusNotFound, err.Error(), h.log)
			return
		}
		httputils.RespondWithError(w, http.StatusInternalServerError, myerrors.ErrInternalServer.Error(), h.log)
		return
	}

	//проверка, что список принадлежит юзеру
	if tasklist.UserID != userID {
		httputils.RespondWithError(w, http.StatusForbidden, myerrors.ErrAccessDenied.Error(), h.log)
		return
	}

	//получение тасок из списка
	tasks, err := h.taskService.GetByListID(r.Context(), taskListID)
	if err != nil {
		httputils.RespondWithError(w, http.StatusInternalServerError, myerrors.ErrInternalServer.Error(), h.log)
		return
	}
	if tasks == nil {
		tasks = []*model.Task{}
	}

	httputils.RespondWithJSON(w, http.StatusOK, tasks, h.log)
}

// /tasks/{id}
func (h *TaskHandlers) GetById(w http.ResponseWriter, r *http.Request) {
	//извлечь user id
	userID, ok := auth.GetUserId(r.Context())
	if !ok {
		httputils.RespondWithError(w, http.StatusUnauthorized, myerrors.ErrUnauthorized.Error(), h.log)
		return
	}

	//излвечь аргумент
	idParam := chi.URLParam(r, "id")
	taskID, err := uuid.Parse(idParam)
	if err != nil {
		httputils.RespondWithError(w, http.StatusNotFound, myerrors.ErrTaskNotFound.Error(), h.log)
		return
	}

	task, err := h.taskService.GetByID(r.Context(), taskID)
	if err != nil {
		if errors.Is(err, myerrors.ErrTaskNotFound) {
			httputils.RespondWithError(w, http.StatusNotFound, err.Error(), h.log)
			return
		}
		httputils.RespondWithError(w, http.StatusInternalServerError, myerrors.ErrInternalServer.Error(), h.log)
		return
	}

	if task.UserID != userID {
		httputils.RespondWithError(w, http.StatusForbidden, myerrors.ErrAccessDenied.Error(), h.log)
		return
	}

	httputils.RespondWithJSON(w, http.StatusOK, task, h.log)
}

// /tasks/{id}
func (h *TaskHandlers) Update(w http.ResponseWriter, r *http.Request) {
	//извлечь user id
	userID, ok := auth.GetUserId(r.Context())
	if !ok {
		httputils.RespondWithError(w, http.StatusUnauthorized, myerrors.ErrUnauthorized.Error(), h.log)
		return
	}

	//собираем новые данные
	updateData := dtoTask.TaskUpdateRequest{}
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
	taskID, err := uuid.Parse(idParam)
	if err != nil {
		httputils.RespondWithError(w, http.StatusNotFound, myerrors.ErrTaskNotFound.Error(), h.log)
		return
	}

	//вызвать метод из сервиса
	updatedTask, err := h.taskService.Update(r.Context(), updateData, taskID, userID)
	if err != nil {
		if errors.Is(err, myerrors.ErrTaskNotFound) {
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
	httputils.RespondWithJSON(w, http.StatusOK, updatedTask, h.log)
}

// /tasks/{id}/status
func (h *TaskHandlers) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	//извлечь user id
	userID, ok := auth.GetUserId(r.Context())
	if !ok {
		httputils.RespondWithError(w, http.StatusUnauthorized, myerrors.ErrUnauthorized.Error(), h.log)
		return
	}

	//собираем новые данные
	updateStatus := dtoTask.TaskUpdateStatusRequest{}
	if err := httputils.ReadFromJSON(r, &updateStatus); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, myerrors.ErrInvalidBody.Error(), h.log)
		return
	}

	//валидировать данные
	if err := updateStatus.Validate(); err != nil {
		httputils.RespondWithError(w, http.StatusBadRequest, err.Error(), h.log)
		return
	}

	//излвечь аргумент
	idParam := chi.URLParam(r, "id")
	taskID, err := uuid.Parse(idParam)
	if err != nil {
		httputils.RespondWithError(w, http.StatusNotFound, myerrors.ErrTaskNotFound.Error(), h.log)
		return
	}

	//вызвать метод из сервиса
	updatedTask, err := h.taskService.UpdateStatus(r.Context(), updateStatus, taskID, userID)
	if err != nil {
		if errors.Is(err, myerrors.ErrTaskNotFound) {
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
	httputils.RespondWithJSON(w, http.StatusOK, updatedTask, h.log)
}

// /tasks/
func (h *TaskHandlers) Create(w http.ResponseWriter, r *http.Request) {
	//извлечь user id
	userID, ok := auth.GetUserId(r.Context())
	if !ok {
		httputils.RespondWithError(w, http.StatusUnauthorized, myerrors.ErrUnauthorized.Error(), h.log)
		return
	}

	//собираем новые данные
	createData := dtoTask.TaskRequest{}
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
	createdTask, err := h.taskService.Create(r.Context(), createData)
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

	//вернуть созданные данные
	w.Header().Set("Location", fmt.Sprintf("/api/tasks/%s", createdTask.ID.String()))
	httputils.RespondWithJSON(w, http.StatusCreated, createdTask, h.log)
}

// /tasks/{id}
func (h *TaskHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	//извлечь user id
	userID, ok := auth.GetUserId(r.Context())
	if !ok {
		httputils.RespondWithError(w, http.StatusUnauthorized, myerrors.ErrUnauthorized.Error(), h.log)
		return
	}

	//излвечь аргумент
	idParam := chi.URLParam(r, "id")
	taskID, err := uuid.Parse(idParam)
	if err != nil {
		httputils.RespondWithError(w, http.StatusNotFound, myerrors.ErrTaskNotFound.Error(), h.log)
		return
	}

	//вызвать метод из сервиса
	if err := h.taskService.Delete(r.Context(), taskID, userID); err != nil {
		if errors.Is(err, myerrors.ErrTaskNotFound) {
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
	w.WriteHeader(http.StatusNoContent)
}
