MIGRATE=migrate -path=./Backend/internal/database/migrations -database "$(DBURL)"

ifneq ($(wildcard .env),)
	DBURL ?= $(strip $(shell grep -E '^DB_HOST=' .env 2>/dev/null | head -n1 | sed 's/^DB_HOST=//'))
endif

export DBURL


## run/backend: run the backend
.PHONY: run/backend
run/backend:
	cd Backend && air

## run/frontend: run the frontend  
.PHONY: run/frontend
run/frontend:
	cd Frontend && npm run dev

## run/authservice: run the auth service
.PHONY: run/authservice
run/authservice:
	cd AuthService && ./mvnw spring-boot:run

# run/redis	
.PHONY: run/redis
run/redis:
	@if [ "$$(docker ps -q -f name=redis-local)" ]; then \
		echo "Container 'redis-local' is already running."; \
	elif [ "$$(docker ps -aq -f name=redis-local)" ]; then \
		echo "Starting existing container 'redis-local'..."; \
		docker start redis-local; \
	else \
		echo "Creating and starting new container 'redis-local'..."; \
		docker run -d -p 6379:6379 --name redis-local redis:7-alpine; \
	fi

## run: run both backend and frontend concurrently
.PHONY: run
run:
	make -j3 run/redis run/backend run/frontend run/authservice


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