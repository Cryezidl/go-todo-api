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
)

type UserRepository struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewUserRepository(db *sqlx.DB, log *slog.Logger) *UserRepository {
	return &UserRepository{db: db, log: log}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	const op = "repository.postgres.UserRepository.Create"
	r.log.Debug("attempting to create user",
		slog.String("op", op),
		slog.String("username", user.Username),
		slog.String("email", user.Email),
	)

	query := `
	INSERT INTO users (
	email, username, password_hash, user_role,
	is_active, timezone, lang, theme) 
	VALUES (
	:email, :username, :password_hash, :user_role,
	:is_active, :timezone, :lang, :theme)
	Returning id, created_at`

	rows, err := r.db.NamedQueryContext(ctx, query, user)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" { // Unique Violation
				if strings.Contains(pqErr.Message, "email") {
					return myerrors.ErrEmailTaken
				}
				if strings.Contains(pqErr.Message, "username") {
					return myerrors.ErrUsernameTaken
				}
			}
		}
		r.log.Error("failed to execute insert query",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("email", user.Email),
		)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Created_At); err != nil {
			r.log.Error("failed to scan returning values",
				slog.String("op", op),
				slog.String("error", err.Error()),
			)
			return err
		}
	}
	r.log.Debug("user created successfully",
		slog.String("op", op),
		slog.String("user_id", user.ID.String()),
	)
	return nil
}

func (r *UserRepository) FindByUserId(ctx context.Context, id uuid.UUID) (*model.User, error) {
	op := "repository.postgres.UserRepository.FindByUserId"

	r.log.Debug("attempting to get user by id",
		slog.String("op", op),
		slog.String("id", id.String()),
	)

	user := &model.User{}
	query := `SELECT * FROM users WHERE id=$1`

	if err := r.db.GetContext(ctx, user, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.Debug("user not found",
				slog.String("op", op),
				slog.String("id", id.String()),
			)
			return nil, myerrors.ErrUserNotFound
		}

		r.log.Error("failed to get user by id",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("id", id.String()),
		)
		return nil, err
	}

	r.log.Debug("user was found",
		slog.String("op", op),
		slog.String("username", user.Username),
		slog.String("id", id.String()),
	)
	return user, nil
}

func (r *UserRepository) FindByUserEmail(ctx context.Context, email string) (*model.User, error) {
	op := "repository.postgres.UserRepository.FindByUserEmail"

	r.log.Debug("attempting to get user by email",
		slog.String("op", op),
		slog.String("email", email),
	)

	user := &model.User{}
	query := `SELECT * FROM users WHERE email=$1`

	if err := r.db.GetContext(ctx, user, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.Debug("user not found",
				slog.String("op", op),
				slog.String("email", email),
			)
			return nil, myerrors.ErrUserNotFound
		}

		r.log.Error("failed to get user by email",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("email", email),
		)
		return nil, err
	}

	r.log.Debug("user was found",
		slog.String("op", op),
		slog.String("username", user.Username),
		slog.String("email", email),
	)
	return user, nil
}

func (r *UserRepository) FindByUserName(ctx context.Context, username string) (*model.User, error) {
	op := "repository.postgres.UserRepository.FindByUserName"

	r.log.Debug("attempting to get user by username",
		slog.String("op", op),
		slog.String("name", username),
	)

	user := &model.User{}
	query := `SELECT * FROM users WHERE username=$1`

	if err := r.db.GetContext(ctx, user, query, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.Debug("user not found",
				slog.String("op", op),
				slog.String("name", username),
			)
			return nil, myerrors.ErrUserNotFound
		}

		r.log.Error("failed to get user by name",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("name", username),
		)
		return nil, err
	}

	r.log.Debug("user was found",
		slog.String("op", op),
		slog.String("id", user.ID.String()),
		slog.String("name", username),
	)
	return user, nil
}

func (r *UserRepository) FindByUserEmailOrName(ctx context.Context, login string) (*model.User, error) {
	op := "repository.postgres.UserRepository.FindByUserEmailOrName"

	r.log.Debug("attempting to get user by username or email",
		slog.String("op", op),
		slog.String("login way", login),
	)

	user := &model.User{}
	query := `SELECT * FROM users WHERE username=$1 OR email=$1 LIMIT 1`

	if err := r.db.GetContext(ctx, user, query, login); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.log.Debug("user not found",
				slog.String("op", op),
				slog.String("login way", login),
			)
			return nil, myerrors.ErrUserNotFound
		}

		r.log.Error("failed to get user by name or email",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("login way", login),
		)
		return nil, err
	}

	r.log.Debug("user was found",
		slog.String("op", op),
		slog.String("id", user.ID.String()),
		slog.String("login way", login),
	)
	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	op := "repository.postgres.UserRepository.Update"
	r.log.Debug("attempting to update user",
		slog.String("op", op),
		slog.String("id", user.ID.String()),
	)

	query := `
	UPDATE USERS
	SET
		email = :email,
		username = :username, 
		password_hash = :password_hash, 
		user_role = :user_role,
		is_active = :is_active, 
		timezone = :timezone, 
		lang = :lang, 
		theme = :theme, 
		updated_at = NOW()
    WHERE id = :id`

	res, err := r.db.NamedExecContext(ctx, query, user)

	if err != nil {
		r.log.Error("failed to update user",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("id", user.ID.String()),
		)
		return err
	}
	updatedRowsCount, _ := res.RowsAffected()
	r.log.Debug("user was updated",
		slog.String("op", op),
		slog.Int64("rows updated", updatedRowsCount),
		slog.String("id", user.ID.String()),
	)
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	op := "repository.postgres.UserRepository.Delete"

	r.log.Debug("attempting to delete user",
		slog.String("op", op),
		slog.String("id", id.String()),
	)

	query := `DELETE FROM users WHERE id=$1`
	rows, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		r.log.Error("failed to delete user",
			slog.String("op", op),
			slog.String("error", err.Error()),
			slog.String("id", id.String()),
		)
		return err
	}
	rowsCount, _ := rows.RowsAffected()
	if rowsCount == 0 {
		r.log.Error("user not found",
			slog.String("op", op),
			slog.String("id", id.String()),
		)
		return myerrors.ErrUserNotFound
	}
	r.log.Debug("user was deleted",
		slog.String("op", op),
		slog.String("id", id.String()),
	)
	return nil
}
