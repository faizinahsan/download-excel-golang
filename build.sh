#!/bin/bash

# Build and Push Script for Download Excel Project

set -e

echo "=== Download Excel Project Build & Deploy Script ==="

# Variables
APP_NAME="download-excel-app"
DOCKER_USER=${DOCKER_USER:-"your-dockerhub-username"}
DOCKER_TAG=${DOCKER_TAG:-"latest"}

# Function to check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        echo "Error: Docker is not running. Please start Docker and try again."
        exit 1
    fi
}

# Function to login to Docker Hub
docker_login() {
    echo "Logging in to Docker Hub..."
    if [ -z "$DOCKER_PASSWORD" ]; then
        echo "Please enter your Docker Hub password:"
        docker login -u "$DOCKER_USER"
    else
        echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USER" --password-stdin
    fi
}

# Function to build Go application
build_go() {
    echo "Building Go application..."
    go mod tidy
    go build -o app main.go
    echo "✓ Go application built successfully"
}

# Function to run migration
run_migration() {
    echo "Running database migration..."
    go run config/postgresql.go
    echo "✓ Database migration completed"
}

# Function to build Docker image
build_docker() {
    echo "Building Docker image..."
    docker build -t "$APP_NAME:$DOCKER_TAG" .
    echo "✓ Docker image built successfully"
}

# Function to push to Docker Hub
push_docker() {
    echo "Tagging and pushing to Docker Hub..."
    docker tag "$APP_NAME:$DOCKER_TAG" "$DOCKER_USER/$APP_NAME:$DOCKER_TAG"
    docker push "$DOCKER_USER/$APP_NAME:$DOCKER_TAG"
    echo "✓ Docker image pushed to Docker Hub successfully"
}

# Function to cleanup
cleanup() {
    echo "Cleaning up..."
    rm -f app
    echo "✓ Local artifacts cleaned up"
}

# Main execution
main() {
    case "${1:-all}" in
        "build")
            build_go
            ;;
        "migrate")
            run_migration
            ;;
        "docker")
            check_docker
            build_docker
            ;;
        "push")
            check_docker
            docker_login
            build_docker
            push_docker
            ;;
        "all")
            echo "Running full build and deploy process..."
            build_go
            run_migration
            check_docker
            docker_login
            build_docker
            push_docker
            cleanup
            echo "🎉 Build and deploy completed successfully!"
            ;;
        "clean")
            cleanup
            ;;
        *)
            echo "Usage: $0 {build|migrate|docker|push|all|clean}"
            echo "  build   - Build Go application only"
            echo "  migrate - Run database migration only"
            echo "  docker  - Build Docker image only"
            echo "  push    - Build and push to Docker Hub"
            echo "  all     - Run complete build and deploy process"
            echo "  clean   - Clean up local artifacts"
            echo ""
            echo "Environment variables:"
            echo "  DOCKER_USER     - Docker Hub username"
            echo "  DOCKER_PASSWORD - Docker Hub password (optional)"
            echo "  DOCKER_TAG      - Docker image tag (default: latest)"
            exit 1
            ;;
    esac
}

main "$@"
