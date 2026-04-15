package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/Cryezidl/go-todo-api/internal/model"
	"github.com/Cryezidl/go-todo-api/pkg/myerrors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type TaskRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewTaskRepository(db *sqlx.DB, log *slog.Logger) *TaskRepository {
	return &TaskRepository{db: db, log: log}
}

func (r *TaskRepository) Create(ctx context.Context, task *model.Task) error {
	const op = "repository.postgres.TaskRepository.Create"
	r.log.Debug("attempting to create task",
		slog.String("op", op),
		slog.String("title", task.Title),
	)

	query := `
	INSERT INTO tasks (
	list_id, user_id, title, description, status, 
	priority, deadline, remind_at, tags) 
	VALUES (
	:list_id, :user_id, :title, :description, :status, 
	:priority, :deadline, :remind_at, :tags)
	RETURNING id, created_at`

	rows, err := r.db.NamedQueryContext(ctx, query, task)

	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&task.ID, &task.CreatedAt); err != nil {
			return err
		}
	}

	return nil
}

func (r *TaskRepository) FindById(ctx context.Context, ID uuid.UUID) (*model.Task, error) {
	op := "repository.postgres.TaskRepository.FindById"

	r.log.Debug("attempting to get task by id",
		slog.String("op", op),
		slog.String("id", ID.String()),
	)

	task := &model.Task{}
	query := `SELECT * FROM tasks WHERE id=$1`

	if err := r.db.GetContext(ctx, task, query, ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.Debug("task not found",
				slog.String("op", op),
				slog.String("id", ID.String()),
			)
			return nil, myerrors.ErrTaskNotFound
		}

		r.log.Error("failed to get task by id",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("id", ID.String()),
		)
		return nil, err
	}

	r.log.Debug("task was found",
		slog.String("op", op),
		slog.String("title", task.Title),
		slog.String("id", ID.String()),
	)
	return task, nil
}

func (r *TaskRepository) FindByListID(ctx context.Context, taskListID uuid.UUID) ([]*model.Task, error) {
	op := "repository.postgres.TaskRepository.FindByListID"

	r.log.Debug("attempting to get task by list id",
		slog.String("op", op),
		slog.String("listID", taskListID.String()),
	)

	tasks := []*model.Task{}
	query := `SELECT * FROM tasks WHERE list_id=$1`

	err := r.db.SelectContext(ctx, &tasks, query, taskListID)
	if err != nil {
		r.log.Error("failed to get task",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepository) Update(ctx context.Context, task *model.Task) error {
	op := "repository.postgres.TaskRepository.Update"
	r.log.Debug("attempting to update task",
		slog.String("op", op),
		slog.String("id", task.ID.String()),
	)

	query := `
	UPDATE tasks
	SET
		list_id = :list_id,
		title = :title,
		description = :description,  
		status = :status,
		priority = :priority,
		remind_at = :remind_at,
		tags = :tags,
		completed_at = :completed_at,
		updated_at = NOW()
    WHERE id = :id`

	rows, err := r.db.NamedQueryContext(ctx, query, task)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return myerrors.ErrTaskNotFound
	}

	if err := rows.Scan(&task.UpdatedAt); err != nil {
		return err
	}
	return nil
}

func (r *TaskRepository) UpdateStatus(ctx context.Context, ID uuid.UUID, status string, completedAt *time.Time) error {
	op := "repository.postgres.TaskRepository.UpdateStatus"
	r.log.Debug("attempting to update task status",
		slog.String("op", op),
		slog.String("id", ID.String()),
		slog.String("status", status),
	)

	query := `
	UPDATE tasks
	SET
		status = $1,
		completed_at = $2,
		updated_at = NOW()
    WHERE id = $3`

	res, err := r.db.ExecContext(ctx, query, status, completedAt, ID)
	if err != nil {
		return err
	}

	rowsCount, _ := res.RowsAffected()
	if rowsCount == 0 {
		return myerrors.ErrTaskNotFound
	}

	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, ID uuid.UUID) error {
	op := "repository.postgres.TaskRepository.Delete"

	r.log.Debug("attempting to delete task",
		slog.String("op", op),
		slog.String("id", ID.String()),
	)

	query := `DELETE FROM tasks WHERE id=$1`
	rows, err := r.db.ExecContext(ctx, query, ID)

	if err != nil {
		r.log.Error("failed to delete task",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("id", ID.String()),
		)
		return err
	}
	rowsCount, _ := rows.RowsAffected()
	if rowsCount == 0 {
		r.log.Error("task not found",
			slog.String("op", op),
			slog.String("id", ID.String()),
		)
		return myerrors.ErrTaskNotFound
	}

	r.log.Debug("task list was deleted",
		slog.String("op", op),
		slog.String("id", ID.String()),
	)
	return nil
}
