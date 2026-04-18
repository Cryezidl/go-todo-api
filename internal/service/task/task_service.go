package task

import (
	"context"
	"errors"
	"log/slog"
	"time"

	dtoTask "github.com/Cryezidl/go-todo-api/internal/dto/task"
	"github.com/Cryezidl/go-todo-api/internal/model"
	"github.com/Cryezidl/go-todo-api/internal/repository"
	"github.com/Cryezidl/go-todo-api/pkg/myerrors"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type TaskService struct {
	taskRepository     repository.TaskRepository
	taskListRepository repository.TaskListRepository
	log                *slog.Logger
}

func NewTaskService(taskRepository repository.TaskRepository, taskListRepository repository.TaskListRepository, log *slog.Logger) *TaskService {
	return &TaskService{taskRepository: taskRepository, taskListRepository: taskListRepository, log: log}
}

func (s *TaskService) Create(ctx context.Context, req dtoTask.TaskRequest) (model.Task, error) {
	const op = "service.TaskService.Create"

	s.log.Debug("creating task",
		slog.String("op", op),
		slog.String("list_id", req.ListID.String()),
		slog.String("user_id", req.UserID.String()),
		slog.String("title", req.Title),
	)

	// Проверяем существование списка
	taskList, err := s.taskListRepository.FindById(ctx, req.ListID)
	if err != nil {
		if errors.Is(err, myerrors.ErrTaskListNotFound) {
			s.log.Warn("task list not found",
				slog.String("op", op),
				slog.String("list_id", req.ListID.String()),
			)
			return model.Task{}, myerrors.ErrTaskListNotFound
		}
		s.log.Error("failed to find task list",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("list_id", req.ListID.String()),
		)
		return model.Task{}, err
	}

	// Проверяем права доступа
	if taskList.UserID != req.UserID {
		s.log.Warn("access denied to create task in list",
			slog.String("op", op),
			slog.String("list_owner_id", taskList.UserID.String()),
			slog.String("request_user_id", req.UserID.String()),
		)
		return model.Task{}, myerrors.ErrAccessDenied
	}

	task := &model.Task{
		ListID:      req.ListID,
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		Deadline:    req.Deadline,
		RemindAt:    req.RemindAt,
		Tags:        pq.StringArray(req.Tags),
		CompletedAt: nil,
	}

	if err := s.taskRepository.Create(ctx, task); err != nil {
		s.log.Error("failed to create task",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("title", req.Title),
		)
		return model.Task{}, err
	}

	s.log.Info("task created successfully",
		slog.String("op", op),
		slog.String("id", task.ID.String()),
		slog.String("title", task.Title),
	)
	return *task, nil
}

func (s *TaskService) GetByID(ctx context.Context, ID uuid.UUID) (model.Task, error) {
	const op = "service.TaskService.GetByID"

	s.log.Debug("getting task by id",
		slog.String("op", op),
		slog.String("id", ID.String()),
	)

	task, err := s.taskRepository.FindById(ctx, ID)
	if err != nil {
		if errors.Is(err, myerrors.ErrTaskNotFound) {
			s.log.Debug("task not found",
				slog.String("op", op),
				slog.String("id", ID.String()),
			)
			return model.Task{}, err
		}
		s.log.Error("failed to get task",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("id", ID.String()),
		)
		return model.Task{}, err
	}

	s.log.Debug("task retrieved",
		slog.String("op", op),
		slog.String("id", task.ID.String()),
		slog.String("title", task.Title),
		slog.String("status", task.Status),
	)
	return *task, nil
}

func (s *TaskService) GetByListID(ctx context.Context, listID uuid.UUID) ([]*model.Task, error) {
	const op = "service.TaskService.GetByListID"

	s.log.Debug("getting tasks by list id",
		slog.String("op", op),
		slog.String("list_id", listID.String()),
	)

	tasks, err := s.taskRepository.FindByListID(ctx, listID)
	if err != nil {
		s.log.Error("failed to get tasks by list id",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("list_id", listID.String()),
		)
		return nil, err
	}

	s.log.Debug("tasks retrieved",
		slog.String("op", op),
		slog.Int("count", len(tasks)),
		slog.String("list_id", listID.String()),
	)
	return tasks, nil
}

func (s *TaskService) Update(ctx context.Context, req dtoTask.TaskUpdateRequest, taskID, userID uuid.UUID) (model.Task, error) {
	const op = "service.TaskService.Update"

	s.log.Debug("updating task",
		slog.String("op", op),
		slog.String("task_id", taskID.String()),
		slog.String("user_id", userID.String()),
	)

	task, err := s.taskRepository.FindById(ctx, taskID)
	if err != nil {
		s.log.Error("failed to find task for update",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("task_id", taskID.String()),
		)
		return model.Task{}, err
	}

	// Проверка прав доступа
	if task.UserID != userID {
		s.log.Warn("access denied for update",
			slog.String("op", op),
			slog.String("task_owner_id", task.UserID.String()),
			slog.String("request_user_id", userID.String()),
			slog.String("task_id", taskID.String()),
		)
		return model.Task{}, myerrors.ErrAccessDenied
	}

	// Логируем изменения
	if req.Title != nil && *req.Title != task.Title {
		s.log.Debug("updating title",
			slog.String("op", op),
			slog.String("old", task.Title),
			slog.String("new", *req.Title),
		)
		task.Title = *req.Title
	}
	if req.Description != nil && req.Description != &task.Description {
		s.log.Debug("updating description",
			slog.String("op", op),
			slog.String("old", task.Description),
			slog.String("new", *req.Description),
		)
		task.Description = *req.Description
	}
	if req.Status != nil && *req.Status != task.Status {
		s.log.Debug("updating status",
			slog.String("op", op),
			slog.String("old", task.Status),
			slog.String("new", *req.Status),
		)
		task.Status = *req.Status
		if *req.Status == "done" {
			now := time.Now()
			task.CompletedAt = &now
			s.log.Debug("task marked as done, setting completed_at",
				slog.String("op", op),
				slog.String("completed_at", now.String()),
			)
		} else {
			task.CompletedAt = nil
			s.log.Debug("task not done, clearing completed_at", slog.String("op", op))
		}
	}
	if req.Priority != nil && *req.Priority != task.Priority {
		s.log.Debug("updating priority",
			slog.String("op", op),
			slog.Int("old", task.Priority),
			slog.Int("new", *req.Priority),
		)
		task.Priority = *req.Priority
	}
	if req.Deadline != nil && req.Deadline != task.Deadline {
		s.log.Debug("updating deadline",
			slog.String("op", op),
			slog.String("old", task.Deadline.String()),
			slog.String("new", req.Deadline.String()),
		)
		task.Deadline = req.Deadline
	}

	if req.RemindAt != nil && req.RemindAt != task.RemindAt {
		s.log.Debug("updating remind_at",
			slog.String("op", op),
			slog.String("old", task.RemindAt.String()),
			slog.String("new", req.RemindAt.String()),
		)
		task.RemindAt = req.RemindAt
	}
	if req.Tags != nil {
		s.log.Debug("updating tags", slog.String("op", op))
		task.Tags = *req.Tags
	}

	if err := s.taskRepository.Update(ctx, task); err != nil {
		s.log.Error("failed to update task in repository",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("task_id", taskID.String()),
		)
		return model.Task{}, err
	}

	s.log.Info("task updated successfully",
		slog.String("op", op),
		slog.String("id", task.ID.String()),
		slog.String("title", task.Title),
	)
	return *task, nil
}

func (s *TaskService) UpdateStatus(ctx context.Context, req dtoTask.TaskUpdateStatusRequest, taskID, userID uuid.UUID) (model.Task, error) {
	const op = "service.TaskService.UpdateStatus"

	s.log.Debug("updating task status",
		slog.String("op", op),
		slog.String("task_id", taskID.String()),
		slog.String("user_id", userID.String()),
		slog.String("new_status", req.Status),
	)

	task, err := s.taskRepository.FindById(ctx, taskID)
	if err != nil {
		s.log.Error("failed to find task for status update",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("task_id", taskID.String()),
		)
		return model.Task{}, err
	}

	// Проверка прав доступа
	if task.UserID != userID {
		s.log.Warn("access denied for status update",
			slog.String("op", op),
			slog.String("task_owner_id", task.UserID.String()),
			slog.String("request_user_id", userID.String()),
			slog.String("task_id", taskID.String()),
		)
		return model.Task{}, myerrors.ErrAccessDenied
	}

	s.log.Debug("current status",
		slog.String("op", op),
		slog.String("old_status", task.Status),
	)

	var completedAt *time.Time
	if req.Status == "done" {
		now := time.Now()
		completedAt = &now
		task.CompletedAt = &now
		s.log.Debug("task completed, setting completed_at",
			slog.String("op", op),
			slog.String("completed_at", now.String()),
		)
	} else {
		completedAt = nil
		task.CompletedAt = nil
		s.log.Debug("task not completed, clearing completed_at", slog.String("op", op))
	}

	if err := s.taskRepository.UpdateStatus(ctx, taskID, req.Status, completedAt); err != nil {
		s.log.Error("failed to update task status in repository",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("task_id", taskID.String()),
			slog.String("status", req.Status),
		)
		return model.Task{}, err
	}

	task.Status = req.Status
	s.log.Info("task status updated successfully",
		slog.String("op", op),
		slog.String("id", task.ID.String()),
		slog.String("new_status", task.Status),
	)
	return *task, nil
}

func (s *TaskService) Delete(ctx context.Context, taskID, userID uuid.UUID) error {
	const op = "service.TaskService.Delete"

	s.log.Debug("deleting task",
		slog.String("op", op),
		slog.String("task_id", taskID.String()),
		slog.String("user_id", userID.String()),
	)

	task, err := s.taskRepository.FindById(ctx, taskID)
	if err != nil {
		s.log.Error("failed to find task for deletion",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("task_id", taskID.String()),
		)
		return err
	}

	// Проверка прав доступа
	if task.UserID != userID {
		s.log.Warn("access denied for deletion",
			slog.String("op", op),
			slog.String("task_owner_id", task.UserID.String()),
			slog.String("request_user_id", userID.String()),
			slog.String("task_id", taskID.String()),
		)
		return myerrors.ErrAccessDenied
	}

	if err := s.taskRepository.Delete(ctx, taskID); err != nil {
		s.log.Error("failed to delete task from repository",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("task_id", taskID.String()),
		)
		return err
	}

	s.log.Info("task deleted successfully",
		slog.String("op", op),
		slog.String("id", taskID.String()),
	)
	return nil
}
