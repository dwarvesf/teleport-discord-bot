# Makefile for Teleport Discord Bot

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY_NAME=teleport-discord-bot
MAIN_PATH=./cmd/teleport-discord-bot

# Environment
export GO111MODULE=on

# Default target
.PHONY: all
all: test build

# Build the application
.PHONY: build
build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Clean up build artifacts
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)
	rm -f coverage.out

# Lint the code
.PHONY: lint
lint:
	golangci-lint run

# Install dependencies
.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run the application locally
.PHONY: run
run:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)
	./$(BINARY_NAME)

# Development setup
.PHONY: setup
setup:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: digest
digest:
	npx ai-digest --whitespace-removal

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all         - Run tests and build the application"
	@echo "  build       - Build the application"
	@echo "  test        - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean       - Remove build artifacts"
	@echo "  lint        - Run golangci-lint"
	@echo "  deps        - Download and tidy dependencies"
	@echo "  run         - Build and run the application"
	@echo "  setup       - Install development tools"
	@echo "  help        - Show this help message"
