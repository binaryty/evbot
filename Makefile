APP_NAME := evbot
BIN_DIR := bin
MAIN_PKG := ./cmd/main.go
CONFIG ?= config/config.yaml
GO_FILES := $(shell git ls-files '*.go')
GOOS := $(shell go env GOOS)

ifeq ($(GOOS),windows)
	MKDIR_CMD := powershell -NoProfile -Command "New-Item -ItemType Directory -Force '$(BIN_DIR)' | Out-Null"
	RM_CMD := powershell -NoProfile -Command "if (Test-Path '$(BIN_DIR)') { Remove-Item -Recurse -Force '$(BIN_DIR)' }"
	SHELL := pwsh.exe
	.SHELLFLAGS := -NoProfile -Command
else
	MKDIR_CMD := mkdir -p $(BIN_DIR)
	RM_CMD := rm -rf $(BIN_DIR)
	SHELL := /bin/sh
endif

.PHONY: help build run test vet lint fmt fmt-check tidy clean ci

help:
	@echo "Common make targets:"
	@echo "  build      - build the bot binary"
	@echo "  run        - run the bot locally (requires CONFIG)"
	@echo "  test       - run unit tests"
	@echo "  lint       - run gofmt check and go vet"
	@echo "  fmt        - format Go files with gofmt"
	@echo "  tidy       - tidy go.mod/go.sum"
	@echo "  clean      - remove build artifacts"
	@echo "  ci         - run the full CI suite"

build:
	@echo ">> building $(APP_NAME)"
	@$(MKDIR_CMD)
	@go build -o $(BIN_DIR)/$(APP_NAME) $(MAIN_PKG)

run:
	@test -n "$(CONFIG)" || (echo "CONFIG path is required" && exit 1)
	@echo ">> running $(APP_NAME) using $(CONFIG)"
	@CONFIG_PATH=$(CONFIG) go run $(MAIN_PKG)

test:
	@echo ">> running tests"
	@go test ./...

vet:
	@echo ">> running go vet"
	@go vet ./...

fmt:
	@echo ">> formatting source files"
	@gofmt -w $(GO_FILES)

fmt-check:
	@echo ">> checking gofmt formatting"
	@go run ./tools/fmtcheck $(GO_FILES)

lint: fmt-check vet

tidy:
	@echo ">> tidying go.mod/go.sum"
	@go mod tidy

clean:
	@echo ">> cleaning build artifacts"
	@$(RM_CMD)

ci: lint test

