.DEFAULT_GOAL := help
.PHONY: help build test # fix "is up to date"

GIT_SHORT_SHA := $(shell git rev-parse --short HEAD)

help: ## Show this help (default)
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

migrate: ## Migrate the DB to the most recent version available
	@./build/migrate

test: ## Run tests
	@go test ./...

coverage-html: ## Coverage html
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

wire:
	@wire ./internal/app/server
	@wire ./internal/app/migrate
	@wire ./internal/app/scheduler
	@wire ./internal/app/consumer

app := 'server'

build: ## Build application (app=server|scheduler|consumer)
	@wire ./internal/app/$(app)
	@go build -ldflags "-X 'exampleapp/internal/infrastructure/di.Version=$(GIT_SHORT_SHA)'" -o ./build/$(app) ./cmd/$(app)
