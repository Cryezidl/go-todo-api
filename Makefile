# Makefile
.PHONY: up down restart logs ps shell db migrate migrate-all clean help

DOCKER_COMPOSE = docker-compose
DOCKER_COMPOSE_ENV = --env-file .env.docker

help:
	@echo "Available commands:"
	@echo "  make up          - Start all services"
	@echo "  make down        - Stop all services"
	@echo "  make restart     - Restart all services"
	@echo "  make logs        - View logs"
	@echo "  make ps          - Show service status"
	@echo "  make shell       - Enter server container"
	@echo "  make db          - Enter database"
	@echo "  make migrate     - Run all migrations"
	@echo "  make migrate-user     - Run users migration"
	@echo "  make migrate-tasklist - Run task lists migration"
	@echo "  make migrate-task     - Run tasks migration"
	@echo "  make clean       - Remove all containers and volumes"

up:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_ENV) up -d --build

down:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_ENV) down

restart: down up

logs:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_ENV) logs -f

ps:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_ENV) ps

shell:
	docker exec -it todo_api sh

db:
	docker exec -it todoapi-db-1 psql -U postgres -d db_todo_api1

# Запустить все миграции
migrate: migrate-user migrate-tasklist migrate-task
	@echo "All migrations completed"

# Отдельные миграции
migrate-user:
	@echo "Running users migration..."
	docker exec -it todoapi-db-1 psql -U postgres -d db_todo_api1 -f /docker-entrypoint-initdb.d/0000_create_users.up.sql

migrate-tasklist:
	@echo "Running task lists migration..."
	docker exec -it todoapi-db-1 psql -U postgres -d db_todo_api1 -f /docker-entrypoint-initdb.d/0001_create_task_lists.up.sql

migrate-task:
	@echo "Running tasks migration..."
	docker exec -it todoapi-db-1 psql -U postgres -d db_todo_api1 -f /docker-entrypoint-initdb.d/0002_create_tasks.up.sql

# Сброс (удалить все таблицы)
migrate-reset:
	@echo "Dropping all tables..."
	docker exec -it todoapi-db-1 psql -U postgres -d db_todo_api1 -c "DROP TABLE IF EXISTS tasks CASCADE;"
	docker exec -it todoapi-db-1 psql -U postgres -d db_todo_api1 -c "DROP TABLE IF EXISTS task_lists CASCADE;"
	docker exec -it todoapi-db-1 psql -U postgres -d db_todo_api1 -c "DROP TABLE IF EXISTS users CASCADE;"
	@echo "Tables dropped. Run 'make migrate' to recreate them."

# Показать все таблицы
migrate-show:
	docker exec -it todoapi-db-1 psql -U postgres -d db_todo_api1 -c "\dt"

clean:
	$(DOCKER_COMPOSE) $(DOCKER_COMPOSE_ENV) down -v
	docker system prune -f