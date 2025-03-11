# Variables
BINARY_NAME=app
DOCKER_IMAGE_NAME=task-manager
DOCKER_TAG=latest

# Build the Go binary
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -v -o $(BINARY_NAME) main.go

# Build Docker image
docker-build:
	@echo "Building Docker image $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)..."
	@docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) .

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_NAME) server

# Run Docker container
docker-run:
	@echo "Running Docker container $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)..."
	@docker run -p 8080:8080 $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) --name $(DOCKER_IMAGE_NAME)

# Run go test for the project
test:
	@echo "Running go test for the project ${DOCKER_IMAGE_NAME}"
	@go test ./...

# Clean up
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	@docker rmi $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)

.PHONY: build docker-build run docker-run clean