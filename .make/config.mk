# Linux / MacOS variables
BINARY_NAME ?= ./bin/$(APP_NAME)
MAIN_PATH ?= ./cmd
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME ?= $(shell date -u +%Y-%m-%d_%H:%M:%S)

# OS commands
MKDIR_P = mkdir -p bin
RM_RF = rm -rf bin/
RM_F = rm -f coverage.out coverage.html