package main

import (
	"exampleapp/cmd/web/internal"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env.local")
	godotenv.Load(".env")
	app := internal.InitializeApp()
	app.Run()
}
