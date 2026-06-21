# go-todo-api

Backend REST API for managing todo tasks, written in Go. Features JWT authentication, task lists, CRUD operations on tasks, and PostgreSQL persistence.

## Tech Stack

- **Go** 1.25
- **chi** — HTTP router & middleware
- **PostgreSQL** + **sqlx** — persistence
- **JWT** (golang-jwt/jwt) — authentication
- **cleanenv** — config loading from `.env` / environment variables
- **testify** — unit testing
- **Docker** & **docker-compose** — containerized setup

## Architecture

Layered structure: `handler → service → repository`

```
cmd/api/          # entrypoint, server bootstrap, graceful shutdown
internal/
  config/         # env-based configuration
  dto/            # request/response payloads
  handlers/       # HTTP handlers
  middleware/     # auth middleware
  model/          # domain models
  repository/     # repository interfaces + postgres implementation
  router/         # route registration per resource
  service/        # business logic
migrations/       # SQL migrations
pkg/               # shared utilities (hashing, JWT, logger, errors, validation)
```

## Getting Started

### Prerequisites

- Go 1.25+
- Docker & Docker Compose (recommended)
- PostgreSQL 16 (if running without Docker)

### Run with Docker (recommended)

1. Create a `.env` file in the project root:

   ```env
   PORT=8080
   APP_ENV=local

   DB_HOST=db
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=todo_db
   DB_SSLMODE=disable

   JWT_SECRET=your-secret-key
   JWT_EXPIRATION_HOURS=24h
   ```

2. Start the stack:

   ```bash
   docker-compose up --build
   ```

   This spins up the API on `:8080` and a PostgreSQL container with migrations applied automatically on first init.

### Run locally (without Docker)

1. Start a local PostgreSQL instance and create a database.
2. Set the environment variables listed above (with `DB_HOST=localhost`).
3. Run:

   ```bash
   go run ./cmd/api
   ```

## Configuration

| Variable | Required | Default | Description |
|---|---|---|---|
| `PORT` | no | `8080` | HTTP server port |
| `APP_ENV` | no | `local` | Application environment |
| `HTTP_READ_TIMEOUT` | no | `5s` | HTTP read timeout |
| `HTTP_WRITE_TIMEOUT` | no | `10s` | HTTP write timeout |
| `DB_HOST` | **yes** | — | PostgreSQL host |
| `DB_PORT` | no | `5432` | PostgreSQL port |
| `DB_USER` | **yes** | — | PostgreSQL user |
| `DB_PASSWORD` | **yes** | — | PostgreSQL password |
| `DB_NAME` | **yes** | — | PostgreSQL database name |
| `DB_SSLMODE` | no | `disable` | PostgreSQL SSL mode |
| `JWT_SECRET` | **yes** | — | Secret used to sign JWTs |
| `JWT_EXPIRATION_HOURS` | no | `24h` | JWT token lifetime |

## API Reference

All routes are prefixed with `/api/v1`. Routes marked 🔒 require a `Bearer` JWT in the `Authorization` header.

### Auth
| Method | Path | Description |
|---|---|---|
| POST | `/auth/register` | Register a new user |
| POST | `/auth/login` | Log in and receive a JWT |

### Users
| Method | Path | Description |
|---|---|---|
| GET | `/users/{username}` | Get public user info by username |
| GET 🔒 | `/users/me` | Get current user profile |
| PATCH 🔒 | `/users/me` | Update current user profile |
| PATCH 🔒 | `/users/me/email` | Update current user's email |
| PATCH 🔒 | `/users/me/password` | Update current user's password |
| DELETE 🔒 | `/users/me` | Delete current user account |

### Task Lists
| Method | Path | Description |
|---|---|---|
| POST 🔒 | `/task-lists` | Create a task list |
| GET 🔒 | `/task-lists` | Get all task lists for current user |
| GET 🔒 | `/task-lists/{id}` | Get a task list by ID |
| PATCH 🔒 | `/task-lists/{id}` | Update a task list |
| DELETE 🔒 | `/task-lists/{id}` | Delete a task list |

### Tasks
| Method | Path | Description |
|---|---|---|
| POST 🔒 | `/tasks` | Create a task |
| GET 🔒 | `/tasks` | Get all tasks |
| GET 🔒 | `/tasks/{id}` | Get a task by ID |
| PATCH 🔒 | `/tasks/{id}` | Update a task |
| PATCH 🔒 | `/tasks/{id}/status` | Update a task's status |
| DELETE 🔒 | `/tasks/{id}` | Delete a task |

## Testing

```bash
go test ./...
```

Unit tests cover the auth and user service layers using mocked repositories.
