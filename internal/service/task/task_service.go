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
	taskList, err := s.taskListRepository.FindById(ctx, req.ListID)
	if err != nil {
		if errors.Is(err, myerrors.ErrTaskListNotFound) {
			return model.Task{}, myerrors.ErrTaskListNotFound
		}
		return model.Task{}, err
	}

	if taskList.UserID != req.UserID {
		return model.Task{}, myerrors.ErrAccessDenied
	}

	task := &model.Task{
		ListID:      req.ListID,
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		RemindAt:    req.RemindAt,
		Tags:        req.Tags,
		CompletedAt: nil,
	}

	return *task, s.taskRepository.Create(ctx, task)
}

func (s *TaskService) GetByID(ctx context.Context, ID uuid.UUID) (model.Task, error) {
	task, err := s.taskRepository.FindById(ctx, ID)
	if err != nil {
		return model.Task{}, err
	}

	return *task, nil
}

func (s *TaskService) GetByListID(ctx context.Context, listID uuid.UUID) ([]*model.Task, error) {
	tasks, err := s.taskRepository.FindByListID(ctx, listID)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *TaskService) Update(ctx context.Context, req dtoTask.TaskUpdateRequest, taskID, userID uuid.UUID) (model.Task, error) {
	task, err := s.taskRepository.FindById(ctx, taskID)
	if err != nil {
		return model.Task{}, err
	}

	if task.UserID != userID {
		return model.Task{}, myerrors.ErrAccessDenied
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Status != nil {
		task.Status = *req.Status
		if *req.Status == "done" {
			now := time.Now()
			task.CompletedAt = &now
		} else {
			task.CompletedAt = nil
		}
	}
	if req.Priority != nil {
		task.Priority = *req.Priority
	}
	if req.RemindAt != nil {
		task.RemindAt = *req.RemindAt
	}
	if req.Tags != nil {
		task.Tags = *req.Tags
	}

	if err := s.taskRepository.Update(ctx, task); err != nil {
		return model.Task{}, err
	}

	return *task, nil
}

func (s *TaskService) UpdateStatus(ctx context.Context, req dtoTask.TaskUpdateStatusRequest, taskID, userID uuid.UUID) (model.Task, error) {
	task, err := s.taskRepository.FindById(ctx, taskID)
	if err != nil {
		return model.Task{}, err
	}

	if task.UserID != userID {
		return model.Task{}, myerrors.ErrAccessDenied
	}

	var completedAt *time.Time
	if req.Status == "done" {
		now := time.Now()
		completedAt = &now
		task.CompletedAt = &now
	} else {
		completedAt = nil
		task.CompletedAt = nil
	}

	if err := s.taskRepository.UpdateStatus(ctx, taskID, req.Status, completedAt); err != nil {
		return model.Task{}, err
	}

	task.Status = req.Status
	return *task, nil
}

func (s *TaskService) Delete(ctx context.Context, taskID, userID uuid.UUID) error {
	task, err := s.taskRepository.FindById(ctx, taskID)
	if err != nil {
		return err
	}

	if task.UserID != userID {
		return myerrors.ErrAccessDenied
	}
	return s.taskRepository.Delete(ctx, taskID)
}
