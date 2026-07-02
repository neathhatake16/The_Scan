
# Makefile for The Scan

.PHONY: build run test clean tidy install migrate-create migrate-up migrate-down migrate-force migrate-drop migrations

DB_URL := postgres://$(shell grep DB_USER backend/.env | cut -d'=' -f2):$(shell grep DB_PASSWORD backend/.env | cut -d'=' -f2)@$(shell grep DB_HOST backend/.env | cut -d'=' -f2):$(shell grep DB_PORT backend/.env | cut -d'=' -f2)/$(shell grep DB_NAME backend/.env | cut -d'=' -f2)?sslmode=disable

MIGRATIONS_DIR := backend/migrations

# Build the application
build:
	cd backend && go build -o bin/server ./cmd/server

run:
	cd backend && go run ./cmd/server

# Help
help:
	@echo "Available commands:"
	@echo "build - Build the application"
	@echo "run - Run the application"
	@echo "migrate-create NAME=<name> - Create a new migration"
	@echo "migrate-up - Run pending migrations"
	@echo "migrate-down - Rollback last migration"
	@echo "migrate-force VERSION=<n> - Force migrate to version"
	@echo "migrate-drop - Drop all tables"

# Create a new migration
migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $$name

# Run all pending migrations
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

# Rollback last migration
migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

# Force migrate to specific version
migrate-force:
	@read -p "Version number: " ver; \
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $$ver

# Drop all tables
migrate-drop:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" drop

# Show migration status
migrations:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

# Go module tidy
tidy:
	cd backend && go mod tidy
