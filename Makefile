include .env

ifeq ($(POSTGRES_SETUP_TEST),)
    POSTGRES_SETUP_TEST := user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) host=$(DB_ADDR) port=$(DB_PORT) sslmode=disable
endif


PROJECT_ROOT := $(CURDIR)
MIGRATION_FOLDER := $(PROJECT_ROOT)/migrations

.PHONY: migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(name)" sql

.PHONY: migration-up
migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

.PHONY: migration-down
migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down

build:
	go build -o bin/server ./cmd/server

go-run-server:
	go run ./cmd/server/main.go

start-services:
	docker-compose up --force-recreate --build -d

run-all: build create-network start-services
	echo "Waiting for containers to start..." && sleep 3
	make migration-up
	docker-compose logs -f

create-network:
# Если сеть уже есть - он её просто выведет, если нет - то создаст
	docker network inspect mynetwork || \
    docker network create mynetwork