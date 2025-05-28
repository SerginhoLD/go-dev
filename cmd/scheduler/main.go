package main

import (
	"exampleapp/cmd/scheduler/internal"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env.local")
	godotenv.Load(".env")
	app := internal.InitializeScheduler()
	app.Run()
}
