# Go application Makefile
APP_NAME ?= account

# Config manager: include OS specific variables
ifeq ($(OS),Windows_NT)
	include .make\config_windows.mk
else
	include ./.make/config.mk
endif

LDFLAGS ?= -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Default target
.DEFAULT_GOAL := build

# Build the application
build: clean
	@echo "Building $(APP_NAME)..."
	@echo $(MKDIR_P)
	go build $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_PATH)

# Run the application
run: build
	@echo "Running $(APP_NAME)..."
	$(BINARY_NAME)

# Run without building (if binary exists)
run-only:
	@echo "Running $(APP_NAME)..."
	$(BINARY_NAME)

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Update dependencies
deps-update:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# Format code
fmt:
	@echo "Formatting code..."
	gofmt -s -w .
	gofumpt -l -w -s .
	goimports -l -w .

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run --timeout=10m -v ./...

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	go test -bench=. ./...

# Clear build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@echo $(RM_RF)
	@echo $(RM_F)

# Install application
install: build
	@echo "Installing $(APP_NAME)"
	go install $(LDFLAGS) $(MAIN_PATH)

# Generate documentation
docs:
	@echo "Generating documentation..."
	godoc -http=:6060

# Show help
help:
	@echo "Available targets:"
	@echo "  build			- Build the application"
	@echo "  run			- Build and run the application"
	@echo "  run-only		- Run the application (without building)"
	@echo "  deps			- Install dependencies"
	@echo "  deps-update	- Update dependencies"
	@echo "  fmt			- Run linter"
	@echo "  test			- Run tests"
	@echo "  test-coverage	- Run tests with coverage report"
	@echo "  bench			- Run benchmarks"
	@echo "  clean			- Clean build artifacts"
	@echo "  install		- Innstall the application"
	@echo "  docs			- Generate documentation"
	@echo "  help			- Show this help"

# Developnem workflow
dev: fmt lint test build

# CI/CD workflow
ci: deps fmt lint test build

.PHONY: build run run-only deps deps-update fmt lint test test-coverage bench clean install docs help dev ci
