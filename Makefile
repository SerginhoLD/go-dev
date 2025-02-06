migrate: # Migrate the DB to the most recent version available
	@goose up


name := app

create-migration: # Creates new migration file with the current timestamp
	@goose create $(name) sql