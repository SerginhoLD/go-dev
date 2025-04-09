.DEFAULT_GOAL := help
.PHONY: help build test # fix "is up to date"

help: ## Show this help (default)
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

migrate: ## Migrate the DB to the most recent version available
	@goose up

down-migration: ## Roll back the version by 1
	@goose down

dev-db-clean: ## Recreate dev database
	@goose down-to 0
	@goose up

test: ## Run tests
	@go test ./...

coverage-html: ## Coverage html
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html


name := 'app'

create-migration: ## Create new migration (name=app)
	@goose create $(name) sql


app := 'web'

build: ## Build application (app=web|scheduler)
	@wire ./cmd/$(app)
	@go build -o ./build/$(app) ./cmd/$(app)
