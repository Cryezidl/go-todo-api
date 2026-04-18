package tasklist

import (
	"context"
	"errors"
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
	const op = "service.TaskListService.Create"

	s.log.Debug("creating task list",
		slog.String("op", op),
		slog.String("user_id", req.UserID.String()),
		slog.String("title", req.Title),
		slog.Bool("is_private", req.IsPrivate),
	)

	taskList := &model.TaskList{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		IsPrivate:   req.IsPrivate,
	}

	err := s.taskListRepository.Create(ctx, taskList)
	if err != nil {
		s.log.Error("failed to create task list",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("title", req.Title),
		)
		return model.TaskList{}, err
	}

	s.log.Info("task list created successfully",
		slog.String("op", op),
		slog.String("id", taskList.ID.String()),
		slog.String("title", taskList.Title),
	)
	return *taskList, nil
}

func (s *TaskListService) GetByID(ctx context.Context, ID uuid.UUID) (model.TaskList, error) {
	const op = "service.TaskListService.GetByID"

	s.log.Debug("getting task list by id",
		slog.String("op", op),
		slog.String("id", ID.String()),
	)

	taskList, err := s.taskListRepository.FindById(ctx, ID)
	if err != nil {
		if errors.Is(err, myerrors.ErrTaskListNotFound) {
			s.log.Debug("task list not found",
				slog.String("op", op),
				slog.String("id", ID.String()),
			)
			return model.TaskList{}, err
		}
		s.log.Error("failed to get task list",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("id", ID.String()),
		)
		return model.TaskList{}, err
	}

	s.log.Debug("task list retrieved",
		slog.String("op", op),
		slog.String("id", taskList.ID.String()),
		slog.String("title", taskList.Title),
	)
	return *taskList, nil
}

func (s *TaskListService) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.TaskList, error) {
	const op = "service.TaskListService.GetByUserID"

	s.log.Debug("getting task lists by user id",
		slog.String("op", op),
		slog.String("user_id", userID.String()),
	)

	taskLists, err := s.taskListRepository.FindByUserID(ctx, userID)
	if err != nil {
		s.log.Error("failed to get user task lists",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("user_id", userID.String()),
		)
		return nil, err
	}

	s.log.Debug("task lists retrieved",
		slog.String("op", op),
		slog.Int("count", len(taskLists)),
		slog.String("user_id", userID.String()),
	)
	return taskLists, nil
}

func (s *TaskListService) Update(ctx context.Context, req dtoTaskList.TaskListUpdateRequest, taskListID, userID uuid.UUID) (model.TaskList, error) {
	const op = "service.TaskListService.Update"

	s.log.Debug("updating task list",
		slog.String("op", op),
		slog.String("task_list_id", taskListID.String()),
		slog.String("user_id", userID.String()),
	)

	taskList, err := s.taskListRepository.FindById(ctx, taskListID)
	if err != nil {
		s.log.Error("failed to find task list for update",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("task_list_id", taskListID.String()),
		)
		return model.TaskList{}, err
	}

	// Проверка прав доступа
	if taskList.UserID != userID {
		s.log.Warn("access denied for update",
			slog.String("op", op),
			slog.String("task_list_id", taskListID.String()),
			slog.String("owner_id", taskList.UserID.String()),
			slog.String("request_user_id", userID.String()),
		)
		return model.TaskList{}, myerrors.ErrAccessDenied
	}

	// Логируем изменения
	if req.Title != nil && *req.Title != taskList.Title {
		s.log.Debug("updating title",
			slog.String("op", op),
			slog.String("old", taskList.Title),
			slog.String("new", *req.Title),
		)
		taskList.Title = *req.Title
	}
	if req.Description != nil && req.Description != &taskList.Description {
		s.log.Debug("updating description",
			slog.String("op", op),
			slog.String("old", taskList.Description),
			slog.String("new", *req.Description),
		)
		taskList.Description = *req.Description
	}
	if req.IsPrivate != nil && *req.IsPrivate != taskList.IsPrivate {
		s.log.Debug("updating is_private",
			slog.String("op", op),
			slog.Bool("old", taskList.IsPrivate),
			slog.Bool("new", *req.IsPrivate),
		)
		taskList.IsPrivate = *req.IsPrivate
	}

	if err := s.taskListRepository.Update(ctx, taskList); err != nil {
		s.log.Error("failed to update task list in repository",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("task_list_id", taskListID.String()),
		)
		return model.TaskList{}, err
	}

	s.log.Info("task list updated successfully",
		slog.String("op", op),
		slog.String("id", taskList.ID.String()),
		slog.String("title", taskList.Title),
	)
	return *taskList, nil
}

func (s *TaskListService) Delete(ctx context.Context, taskListID, userID uuid.UUID) error {
	const op = "service.TaskListService.Delete"

	s.log.Debug("deleting task list",
		slog.String("op", op),
		slog.String("task_list_id", taskListID.String()),
		slog.String("user_id", userID.String()),
	)

	taskList, err := s.taskListRepository.FindById(ctx, taskListID)
	if err != nil {
		s.log.Error("failed to find task list for deletion",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("task_list_id", taskListID.String()),
		)
		return err
	}

	// Проверка прав доступа
	if taskList.UserID != userID {
		s.log.Warn("access denied for deletion",
			slog.String("op", op),
			slog.String("task_list_id", taskListID.String()),
			slog.String("owner_id", taskList.UserID.String()),
			slog.String("request_user_id", userID.String()),
		)
		return myerrors.ErrAccessDenied
	}

	if err := s.taskListRepository.Delete(ctx, taskListID); err != nil {
		s.log.Error("failed to delete task list from repository",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("task_list_id", taskListID.String()),
		)
		return err
	}

	s.log.Info("task list deleted successfully",
		slog.String("op", op),
		slog.String("id", taskListID.String()),
	)
	return nil
}
