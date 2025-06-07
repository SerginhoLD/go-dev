package main

import (
	"context"
	"exampleapp/internal/app/migrate"
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env.local")
	godotenv.Load(".env")

	conn := migrate.InitializeConn()
	defer conn.Close()

	for _, sql := range list() {
		_, err := conn.ExecContext(context.Background(), sql)

		if err != nil {
			println(fmt.Sprintf("Error: %s", err))
			return
		}
	}

	println("Succcess!")
}

func list() []string {
	return []string{
		`CREATE TABLE IF NOT EXISTS objects(title text, metro string, price int, checked bool, size float, rooms int, updated_at timestamp) morphology = 'lemmatize_ru_all'`,
		`REPLACE INTO objects (id, title, metro, price, size, rooms, updated_at) VALUES(1, '1 комн, 33 м², 2/5 эт', 'Одинцово', 0, 33, 1, '2025-01-01')`,
		`REPLACE INTO objects (id, title, metro, price, size, rooms, updated_at) VALUES(2, '2 комн, 37 м², 1/9 эт', 'Одинцово', 0, 37, 2, '2025-01-01')`,
	}
}
