package main

import (
	"exampleapp/internal/app/scheduler"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env.local")
	godotenv.Load(".env")
	app := scheduler.Initialize()
	app.Run()
}
