.DEFAULT_GOAL := help

help: ## Show this help (default)
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

migrate: ## Migrate the DB to the most recent version available
	@goose up

down-migration: ## Roll back the version by 1
	@goose down


name := app

create-migration: ## Creates new migration file with the current timestamp
	@goose create $(name) sql
