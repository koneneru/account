# Windows CMD variables
BINARY_NAME = bin\$(APP_NAME).exe
MAIN_PATH = ./cmd
VERSION ?= $(shell git describe --tags --always --dirty 2>nul || echo "dev")
BUILD_TIME ?= $(shell powershell -Command "Get-Date -uformat '%Y-%m-%d_%H:%M:%S'")

# OS commands
MKDIR_P = if not exist bin mkdir bin
RM_RF = if exist bin rmdir /s /q bin
RM_F = if exist coverage.out del /f /q coverage.out && if exist coverage.html del /f /q coverage.html