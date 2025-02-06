migrate: # Migrate the DB to the most recent version available
	@goose up

down-migration: # Roll back the version by 1
	@goose down


name := app

create-migration: # Creates new migration file with the current timestamp
	@goose create $(name) sql