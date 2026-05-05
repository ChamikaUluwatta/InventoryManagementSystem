MIGRATE=migrate -path=./Backend/internal/database/migrations -database "$(DB_URL)"

ifneq ($(wildcard .env),)
	DB_URL ?= $(strip $(shell grep -E '^DB_HOST=' .env 2>/dev/null | head -n1 | sed 's/^DB_HOST=//'))
endif

export DB_URL

# Development

## run/backend: run the backend
.PHONY: run/backend
run/backend:
	cd Backend && go run ./cmd

## run/frontend: run the frontend  
.PHONY: run/frontend
run/frontend:
	cd Frontend && npm run dev

## run: run both backend and frontend concurrently
.PHONY: run
run:
	make -j2 run/backend run/frontend

# Database Migrations

## migrate/up: apply all migrations
.PHONY: migrate/up
migrate/up:
	$(MIGRATE) up

## migrate/up/1: apply 1 migration step
.PHONY: migrate/up/1
migrate/up/1:
	$(MIGRATE) up 1

## migrate/down: rollback all migrations
.PHONY: migrate/down
migrate/down:
	$(MIGRATE) down

## migrate/down/1: rollback 1 migration step
.PHONY: migrate/down/1
migrate/down/1:
	$(MIGRATE) down 1

## migrate/create: create a new migration file
.PHONY: migrate/create
migrate/create:
	migrate create -ext sql -seq -dir ./Backend/internal/database/migrations $(name)

# Docker
## docker/up: start with dockerized DB
.PHONY: docker/up/dockerDB
docker/up:
	docker compose --profile dockerDB up

## docker/down: stop all containers
.PHONY: docker/down/dockerDB
docker/down:
	docker compose --profile dockerDB down

.PHONY: docker/up/localDB
docker/up/localDB:
	docker compose --profile localDB up

.PHONY: docker/down/localDB
docker/down/localDB:
	docker compose --profile localDB down


# Help
## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'