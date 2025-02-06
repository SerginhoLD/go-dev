package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"os"
)

func NewDB() *sql.DB {
	dsn, _ := os.LookupEnv("GOOSE_DBSTRING")
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		panic(err)
	}

	//defer db.Close()

	return db
}
