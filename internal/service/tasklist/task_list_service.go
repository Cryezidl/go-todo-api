package task_list

import (
	"context"
	"log/slog"

	dtoTaskList "github.com/Cryezidl/go-todo-api/internal/dto/tasklist"
	"github.com/Cryezidl/go-todo-api/internal/model"
	"github.com/Cryezidl/go-todo-api/internal/repository"
	"github.com/google/uuid"
)

type TaskListService struct {
	TaskListRepository repository.TaskListRepository
	log                *slog.Logger
}

func NewTaskListService(TaskListRepository repository.TaskListRepository, log *slog.Logger) *TaskListService {
	return &TaskListService{TaskListRepository: TaskListRepository, log: log}
}

func (s *TaskListService) Create(ctx context.Context, req dtoTaskList.TaskListRequest) (model.TaskList, error) {
	taskList := &model.TaskList{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		IsPrivate:   req.IsPrivate,
	}

	err := s.TaskListRepository.Create(ctx, taskList)
	return *taskList, err
}

func (s *TaskListService) GetByID(ctx context.Context, ID uuid.UUID) (model.TaskList, error) {
	taskList, err := s.TaskListRepository.FindById(ctx, ID)
	if err != nil {
		return model.TaskList{}, err
	}

	return *taskList, nil
}

func (s *TaskListService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.TaskList, error) {
	taskList, err := s.TaskListRepository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return taskList, nil
}

func (s *TaskListService) Update(ctx context.Context, req dtoTaskList.TaskListUpdateRequest, taskListID uuid.UUID) (model.TaskList, error) {
	taskList, err := s.TaskListRepository.FindById(ctx, taskListID)
	if err != nil {
		return model.TaskList{}, err
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

	if err := s.TaskListRepository.Update(ctx, taskList); err != nil {
		return model.TaskList{}, err
	}

	return *taskList, nil
}

func (s *TaskListService) Delete(ctx context.Context, taskListID uuid.UUID) error {
	return s.TaskListRepository.Delete(ctx, taskListID)
}
