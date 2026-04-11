package tasklist

import (
	"context"
	"log/slog"

	dtoTaskList "github.com/Cryezidl/go-todo-api/internal/dto/tasklist"
	"github.com/Cryezidl/go-todo-api/internal/model"
	"github.com/Cryezidl/go-todo-api/internal/repository"
	"github.com/Cryezidl/go-todo-api/pkg/myerrors"
	"github.com/google/uuid"
)

type TaskListService struct {
	taskListRepository repository.TaskListRepository
	log                *slog.Logger
}

func NewTaskListService(taskListRepository repository.TaskListRepository, log *slog.Logger) *TaskListService {
	return &TaskListService{taskListRepository: taskListRepository, log: log}
}

func (s *TaskListService) Create(ctx context.Context, req dtoTaskList.TaskListRequest) (model.TaskList, error) {
	taskList := &model.TaskList{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		IsPrivate:   req.IsPrivate,
	}

	err := s.taskListRepository.Create(ctx, taskList)
	return *taskList, err
}

func (s *TaskListService) GetByID(ctx context.Context, ID uuid.UUID) (model.TaskList, error) {
	taskList, err := s.taskListRepository.FindById(ctx, ID)
	if err != nil {
		return model.TaskList{}, err
	}

	return *taskList, nil
}

func (s *TaskListService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.TaskList, error) {
	taskList, err := s.taskListRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return taskList, nil
}

func (s *TaskListService) Update(ctx context.Context, req dtoTaskList.TaskListUpdateRequest, taskListID, userID uuid.UUID) (model.TaskList, error) {
	taskList, err := s.taskListRepository.FindById(ctx, taskListID)
	if err != nil {
		return model.TaskList{}, err
	}

	if taskList.UserID != userID {
		return model.TaskList{}, myerrors.ErrAccessDenied
	}

	if req.Title != nil {
		taskList.Title = *req.Title
	}
	if req.Description != nil {
		taskList.Description = *req.Description
	}
	if req.IsPrivate != nil {
		taskList.IsPrivate = *req.IsPrivate
	}

	if err := s.taskListRepository.Update(ctx, taskList); err != nil {
		return model.TaskList{}, err
	}

	return *taskList, nil
}

func (s *TaskListService) Delete(ctx context.Context, taskListID, userID uuid.UUID) error {
	taskList, err := s.taskListRepository.FindById(ctx, taskListID)
	if err != nil {
		return err
	}

	if taskList.UserID != userID {
		return myerrors.ErrAccessDenied
	}
	return s.taskListRepository.Delete(ctx, taskListID)
}
