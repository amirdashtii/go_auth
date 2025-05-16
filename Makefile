.PHONY: build run test clean migrate-up migrate-down mock docker-up docker-down dev run-dev

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=go_auth
MAIN_FILE=cmd/main.go
TMP_DIR=tmp

# Database parameters
DB_HOST=localhost
DB_PORT=5432
DB_USER=go_auth
DB_PASSWORD=go_auth
DB_NAME=go_auth
MIGRATIONS_PATH=migrations

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_FILE)

# Run the application
run:
	$(GOCMD) run $(MAIN_FILE)

# Run in development mode with hot reload
dev: migrate-up run-dev

# Run with hot reload using CompileDaemon
run-dev:
	CompileDaemon --build="go build -o $(TMP_DIR)/main $(MAIN_FILE)" --command="./$(TMP_DIR)/main"

# Run tests
test:
	$(GOTEST) -v ./...

# Run integration tests
test-integration:
	$(GOTEST) -v ./tests/integration/...

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(TMP_DIR)

# Download dependencies
deps:
	$(GOMOD) download

# Generate mocks
mock:
	mockery --all

# Database migrations
migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

# Docker commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# Development setup
setup: deps mock migrate-up

# Help command
help:
	@echo "Available commands:"
	@echo "  build            - Build the application"
	@echo "  run              - Run the application"
	@echo "  dev              - Run in development mode with hot reload"
	@echo "  run-dev          - Run with hot reload using CompileDaemon"
	@echo "  test             - Run all tests"
	@echo "  test-integration - Run integration tests"
	@echo "  clean            - Clean build files"
	@echo "  deps             - Download dependencies"
	@echo "  mock             - Generate mocks"
	@echo "  migrate-up       - Run database migrations"
	@echo "  migrate-down     - Rollback database migrations"
	@echo "  docker-up        - Start Docker containers"
	@echo "  docker-down      - Stop Docker containers"
	@echo "  setup            - Setup development environment" 