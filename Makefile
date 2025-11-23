.PHONY: up down build test test-integration migrate-up migrate-down logs clean

up:
	docker compose up --build -d

down:
	docker compose down

build:
	go build -o bin/server ./cmd/server

test:
	go test -v ./...

test-integration:
	go test -v ./tests/...

migrate-up:
	docker compose exec backend goose -dir migrations postgres "host=db port=5432 user=${POSTGRES_USER} password=$$(cat secrets/postgres-password.txt) dbname=${POSTGRES_DB} sslmode=disable" up

migrate-down:
	docker compose exec backend goose -dir migrations postgres "host=db port=5432 user=${POSTGRES_USER} password=$$(cat secrets/postgres-password.txt) dbname=${POSTGRES_DB} sslmode=disable" down

logs:
	docker compose logs -f backend

clean:
	docker compose down -v
