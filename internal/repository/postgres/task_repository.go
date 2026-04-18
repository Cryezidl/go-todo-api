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
		slog.String("list_id", task.ListID.String()),
		slog.String("user_id", task.UserID.String()),
	)

	query := `
	INSERT INTO tasks (
		list_id, user_id, title, description, status, 
		priority, deadline, remind_at, tags
	) VALUES (
		:list_id, :user_id, :title, :description, :status, 
		:priority, :deadline, :remind_at, :tags
	)
	RETURNING id, created_at, updated_at`

	rows, err := r.db.NamedQueryContext(ctx, query, task)
	if err != nil {
		r.log.Error("failed to execute insert",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("title", task.Title),
		)
		return err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt); err != nil {
			r.log.Error("failed to scan returned values",
				slog.String("op", op),
				slog.String("error", err.Error()),
			)
			return err
		}
		r.log.Debug("task created successfully",
			slog.String("op", op),
			slog.String("id", task.ID.String()),
			slog.String("title", task.Title),
		)
	} else {
		r.log.Error("no rows returned after insert", slog.String("op", op))
		return errors.New("no rows returned after insert")
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
	query := `SELECT * FROM tasks WHERE id = $1`

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
		slog.String("id", task.ID.String()),
		slog.String("title", task.Title),
		slog.String("status", task.Status),
	)
	return task, nil
}

func (r *TaskRepository) FindByListID(ctx context.Context, taskListID uuid.UUID) ([]*model.Task, error) {
	op := "repository.postgres.TaskRepository.FindByListID"

	r.log.Debug("attempting to get tasks by list id",
		slog.String("op", op),
		slog.String("list_id", taskListID.String()),
	)

	var tasks []*model.Task
	query := `SELECT * FROM tasks WHERE list_id = $1 ORDER BY created_at DESC`

	err := r.db.SelectContext(ctx, &tasks, query, taskListID)
	if err != nil {
		r.log.Error("failed to get tasks by list id",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("list_id", taskListID.String()),
		)
		return nil, err
	}

	r.log.Debug("tasks retrieved successfully",
		slog.String("op", op),
		slog.Int("count", len(tasks)),
		slog.String("list_id", taskListID.String()),
	)
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
		deadline = :deadline,
		remind_at = :remind_at,
		tags = :tags,
		completed_at = :completed_at,
		updated_at = NOW()
	WHERE id = :id
	RETURNING updated_at`

	rows, err := r.db.NamedQueryContext(ctx, query, task)
	if err != nil {
		r.log.Error("failed to execute update",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("id", task.ID.String()),
		)
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		r.log.Warn("task not found for update",
			slog.String("op", op),
			slog.String("id", task.ID.String()),
		)
		return myerrors.ErrTaskNotFound
	}

	if err := rows.Scan(&task.UpdatedAt); err != nil {
		r.log.Error("failed to scan updated_at",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("id", task.ID.String()),
		)
		return err
	}

	r.log.Debug("task updated successfully",
		slog.String("op", op),
		slog.String("id", task.ID.String()),
		slog.String("updated_at", task.UpdatedAt.String()),
	)
	return nil
}

func (r *TaskRepository) UpdateStatus(ctx context.Context, ID uuid.UUID, status string, completedAt *time.Time) error {
	op := "repository.postgres.TaskRepository.UpdateStatus"
	r.log.Debug("attempting to update task status",
		slog.String("op", op),
		slog.String("id", ID.String()),
		slog.String("status", status),
	)

	if completedAt != nil {
		r.log.Debug("setting completed_at",
			slog.String("op", op),
			slog.String("completed_at", completedAt.String()),
		)
	} else {
		r.log.Debug("clearing completed_at", slog.String("op", op))
	}

	query := `
	UPDATE tasks
	SET
		status = $1,
		completed_at = $2,
		updated_at = NOW()
	WHERE id = $3`

	res, err := r.db.ExecContext(ctx, query, status, completedAt, ID)
	if err != nil {
		r.log.Error("failed to update task status",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("id", ID.String()),
			slog.String("status", status),
		)
		return err
	}

	rowsCount, err := res.RowsAffected()
	if err != nil {
		r.log.Error("failed to get rows affected",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		return err
	}

	if rowsCount == 0 {
		r.log.Warn("task not found for status update",
			slog.String("op", op),
			slog.String("id", ID.String()),
		)
		return myerrors.ErrTaskNotFound
	}

	r.log.Debug("task status updated successfully",
		slog.String("op", op),
		slog.String("id", ID.String()),
		slog.String("status", status),
		slog.Int64("rows_affected", rowsCount),
	)
	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, ID uuid.UUID) error {
	op := "repository.postgres.TaskRepository.Delete"

	r.log.Debug("attempting to delete task",
		slog.String("op", op),
		slog.String("id", ID.String()),
	)

	query := `DELETE FROM tasks WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, ID)
	if err != nil {
		r.log.Error("failed to delete task",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("id", ID.String()),
		)
		return err
	}

	rowsCount, err := result.RowsAffected()
	if err != nil {
		r.log.Error("failed to get rows affected",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		return err
	}

	if rowsCount == 0 {
		r.log.Warn("task not found for deletion",
			slog.String("op", op),
			slog.String("id", ID.String()),
		)
		return myerrors.ErrTaskNotFound
	}

	r.log.Debug("task deleted successfully",
		slog.String("op", op),
		slog.String("id", ID.String()),
		slog.Int64("rows_affected", rowsCount),
	)
	return nil
}
