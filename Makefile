# Makefile for Download Excel Project

.PHONY: build run docker-build docker-push migrate migrate-create migrate-down clean help

# Variables
APP_NAME=download-excel-app
DOCKER_USER=faizinahsan
DOCKER_TAG=latest

# Build the Go application
build:
	@echo "Building Go application..."
	go build -o app main.go

# Run the application
run:
	@echo "Running the application..."
	./app

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(APP_NAME):$(DOCKER_TAG) .

# Tag and push to Docker Hub
docker-push: docker-build
	@echo "Tagging and pushing to Docker Hub..."
	docker tag $(APP_NAME):$(DOCKER_TAG) $(DOCKER_USER)/$(APP_NAME):$(DOCKER_TAG)
	docker push $(DOCKER_USER)/$(APP_NAME):$(DOCKER_TAG)

# Run database migration with goose
migrate:
	@echo "Running database migration with goose..."
	goose -dir migrations postgres "user=postgres password=password dbname=excel_db sslmode=disable" up

# Create new migration
migrate-create:
	@echo "Creating new migration..."
	@read -p "Enter migration name: " name; \
	goose -dir migrations create $$name sql

# Rollback migration
migrate-down:
	@echo "Rolling back migration..."
	goose -dir migrations postgres "user=postgres password=password dbname=excel_db sslmode=disable" down

# Clean up
clean:
	@echo "Cleaning up..."
	rm -f app
	docker rmi $(APP_NAME):$(DOCKER_TAG) 2>/dev/null || true
	docker rmi $(DOCKER_USER)/$(APP_NAME):$(DOCKER_TAG) 2>/dev/null || true

# Help command
help:
	@echo "Available commands:"
	@echo "  build       - Build the Go application"
	@echo "  run         - Run the application"
	@echo "  docker-build- Build Docker image"
	@echo "  docker-push - Build and push Docker image to Docker Hub"
	@echo "  migrate     - Run database migration with goose up"
	@echo "  migrate-create - Create new migration file"
	@echo "  migrate-down - Rollback last migration"
	@echo "  clean       - Clean up build artifacts and Docker images"
	@echo "  help        - Show this help message"
	@echo ""
	@echo "Note: Make sure goose is installed: go install github.com/pressly/goose/v3/cmd/goose@latest"
