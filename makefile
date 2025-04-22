# Variables
APP_NAME=go-snake
BUILD_DIR=build
MAIN_PATH=cmd/go-snake.go

# Go commands
GO=go
GOBUILD=$(GO) build
GORUN=$(GO) run
GOTEST=$(GO) test
GOMOD=$(GO) mod

# Build tags (Removed noaudio tag as Ebiten handles audio differently)
BUILD_TAGS=

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_TAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

# Run the application
.PHONY: run
run:
	$(GORUN) $(MAIN_PATH)

# Clean up build artifacts
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

# Test the application
.PHONY: test
test:
	$(GOTEST) ./...

# Update dependencies
.PHONY: deps
deps:
	$(GOMOD) tidy
	$(GOMOD) download

# Removed install-deps target as it was specific to g3n/OpenAL

# Run the application with specific config
.PHONY: debug
debug:
	$(GORUN) -race $(MAIN_PATH)

# Help target
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make                - Build the application"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make test           - Run tests"
	@echo "  make deps           - Update dependencies"
	@echo "  make debug          - Run with race detection"
	@echo "  make help           - Show this help"
