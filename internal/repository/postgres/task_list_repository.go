package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strings"

	"github.com/Cryezidl/go-todo-api/internal/model"
	"github.com/Cryezidl/go-todo-api/pkg/myerrors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type TaskListRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewTaskListRepository(db *sqlx.DB, log *slog.Logger) *TaskListRepository {
	return &TaskListRepository{db: db, log: log}
}

func (r *TaskListRepository) Create(ctx context.Context, taskList *model.TaskList) error {
	const op = "repository.postgres.TaskListRepository.Create"
	r.log.Debug("attempting to create task list",
		slog.String("op", op),
		slog.String("title", taskList.Title),
	)

	query := `
	INSERT INTO task_lists (
	user_id, title, description, is_private) 
	VALUES (
	:user_id, :title, :description, :is_private)
	RETURNING id, created_at`

	rows, err := r.db.NamedQueryContext(ctx, query, taskList)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" { // Unique Violation
				if strings.Contains(pqErr.Message, "title") {
					return myerrors.ErrTaskListAlreadyExists
				}
				return err
			}
		}
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&taskList.CreatedAt, &taskList.UpdatedAt); err != nil {
			return err
		}
	}

	return nil
}

func (r *TaskListRepository) FindById(ctx context.Context, ID uuid.UUID) (*model.TaskList, error) {
	op := "repository.postgres.TaskListRepository.FindById"

	r.log.Debug("attempting to get task list by id",
		slog.String("op", op),
		slog.String("id", ID.String()),
	)

	taskList := &model.TaskList{}
	query := `SELECT * FROM task_lists WHERE id=$1`

	if err := r.db.GetContext(ctx, taskList, query, ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.Debug("task list not found",
				slog.String("op", op),
				slog.String("id", ID.String()),
			)
			return nil, myerrors.ErrTaskListNotFound
		}

		r.log.Error("failed to get task lsit by id",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("id", ID.String()),
		)
		return nil, err
	}

	r.log.Debug("task list was found",
		slog.String("op", op),
		slog.String("title", taskList.Title),
		slog.String("id", ID.String()),
	)
	return taskList, nil
}

func (r *TaskListRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*model.TaskList, error) {
	op := "repository.postgres.TaskListRepository.FindByUserID"

	r.log.Debug("attempting to get task list by user id",
		slog.String("op", op),
		slog.String("userID", userID.String()),
	)

	taskLists := []*model.TaskList{}
	query := `SELECT * FROM task_lists WHERE user_id=$1`

	err := r.db.SelectContext(ctx, &taskLists, query, userID)
	if err != nil {
		r.log.Error("failed to get task lists",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	return taskLists, nil
}

func (r *TaskListRepository) Update(ctx context.Context, taskList *model.TaskList) error {
	op := "repository.postgres.TaskListRepository.Update"
	r.log.Debug("attempting to update task list",
		slog.String("op", op),
		slog.String("id", taskList.ID.String()),
	)

	query := `
	UPDATE task_lists
	SET
		title = :title,
		description = :description,  
		is_private = :is_private,
		updated_at = NOW()
    WHERE id = :id
	RETURNING`

	rows, err := r.db.NamedQueryContext(ctx, query, taskList)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return myerrors.ErrTaskListNotFound
	}

	if err := rows.Scan(&taskList.UpdatedAt); err != nil {
		return err
	}
	return nil
}
func (r *TaskListRepository) Delete(ctx context.Context, ID uuid.UUID) error {
	op := "repository.postgres.TaskListRepository.Delete"

	r.log.Debug("attempting to delete task list",
		slog.String("op", op),
		slog.String("id", ID.String()),
	)

	query := `DELETE FROM task_lists WHERE id=$1`
	rows, err := r.db.ExecContext(ctx, query, ID)

	if err != nil {
		r.log.Error("failed to delete task list",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("id", ID.String()),
		)
		return err
	}
	rowsCount, _ := rows.RowsAffected()
	if rowsCount == 0 {
		r.log.Error("task list not found",
			slog.String("op", op),
			slog.String("id", ID.String()),
		)
		return myerrors.ErrTaskListNotFound
	}

	r.log.Debug("task list was deleted",
		slog.String("op", op),
		slog.String("id", ID.String()),
	)
	return nil
}
