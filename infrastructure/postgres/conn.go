package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func NewDB() *sql.DB {
	dsnStr := "user=postgres password=root dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", dsnStr)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	return db
}
