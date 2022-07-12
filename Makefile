#! /usr/bin/make -f

# Project variables.
PROJECT_NAME = RangeProcessor
FIND_ARGS := -name '*.go' -type f -not -name '*.pb.go'
BUILD_FOLDER = ./dist
GOCILINT := $(GOPATH)/bin/golangci-lint
GOIMPORTS := $(GOPATH)/bin/goimports

## install: Install de binary.
install:
	@echo Installing RangeProcessor...
	@go install ./...

## build: Build the binary.
build:
	@echo Building RangeProcessor...
	@-mkdir -p $(BUILD_FOLDER) 2> /dev/null
	@go build -o $(BUILD_FOLDER) ./...

## clean: Clean build files. Also runs `go clean` internally.
clean:
	@echo Cleaning build cache...
	@-rm -rf $(BUILD_FOLDER) 2> /dev/null
	@go clean ./...

## govet: Run go vet.
govet:
	@echo Running go vet...
	@go vet ./...

$(GOIMPORTS):
	@echo Installing goimports...
	@cd && go install golang.org/x/tools/cmd/goimports@latest

## format: Run gofmt.
format:
	@echo Formatting...
	@find . $(FIND_ARGS) | xargs gofmt -d -s
	@find . $(FIND_ARGS) | xargs $(GOIMPORTS) -w -local github.com/dapper-labs-talent/cc_DaniloPantani_Flow-RangeProcessor

$(GOCILINT):
	@echo Installing gocilint...
	@cd && go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.45.2

## lint: Run Golang CI Lint.
lint:
	@echo Running gocilint...
	@$(GOCILINT) run --out-format=tab --issues-exit-code=0

## test-unit: Run the unit tests.
test-unit:
	@echo Running unit tests...
	@go test -race -failfast -v ./...

## test-integration: Run the integration tests.
test-integration: install
	@echo Running integration tests...
	@go test -race -failfast -v -timeout 60m ./integration/...

## test: Run unit and integration tests.
test: govet test-unit test-integration

help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECT_NAME)", or just run 'make' for install"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

.DEFAULT_GOAL := install
