.DEFAULT_GOAL := help

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


name := app

create-migration: ## Creates new migration file with the current timestamp
	@goose create $(name) sql
