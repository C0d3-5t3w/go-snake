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

# Build tags
# Use -tags=noaudio to build without audio support
BUILD_TAGS=

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_TAGS) -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

# Build without audio support (for systems missing OpenAL/Vorbis libraries)
.PHONY: build-noaudio
build-noaudio:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -tags=noaudio -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

# Run the application
.PHONY: run
run:
	$(GORUN) $(MAIN_PATH)

# Run without audio
.PHONY: run-noaudio
run-noaudio:
	$(GORUN) -tags=noaudio $(MAIN_PATH)

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

# Install system dependencies based on platform
.PHONY: install-deps
install-deps:
	@echo "Installing system dependencies for G3N engine..."
	@case "$$(uname -s)" in \
		Linux*) \
			echo "Detected Linux system" && \
			echo "Installing OpenAL and Vorbis libraries..." && \
			(command -v apt-get >/dev/null && sudo apt-get install -y libopenal-dev libvorbis-dev) || \
			(command -v dnf >/dev/null && sudo dnf install -y openal-soft-devel libvorbis-devel) || \
			echo "Please install OpenAL and Vorbis development libraries manually" ;; \
		Darwin*) \
			echo "Detected macOS system" && \
			echo "Installing OpenAL and Vorbis libraries using Homebrew..." && \
			(command -v brew >/dev/null && brew install openal-soft libvorbis) || \
			echo "Please install Homebrew and then run: brew install openal-soft libvorbis" ;; \
		MINGW*|MSYS*) \
			echo "Detected Windows system" && \
			echo "Please install OpenAL and Vorbis development libraries manually" ;; \
		*) \
			echo "Unknown operating system. Please install OpenAL and Vorbis development libraries manually" ;; \
	esac

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
	@echo "  make build-noaudio  - Build without audio support (for systems missing OpenAL/Vorbis)"
	@echo "  make run            - Run the application"
	@echo "  make run-noaudio    - Run without audio support"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make test           - Run tests"
	@echo "  make deps           - Update dependencies"
	@echo "  make install-deps   - Install system dependencies (OpenAL, Vorbis)"
	@echo "  make debug          - Run with race detection"
	@echo "  make help           - Show this help"
