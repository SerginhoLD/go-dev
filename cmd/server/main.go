package main

import (
	"exampleapp/internal/app/server"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env.local")
	godotenv.Load(".env")
	app := server.Initialize()
	app.Run()
}
