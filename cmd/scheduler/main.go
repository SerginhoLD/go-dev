package main

import (
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env.local")
	godotenv.Load(".env")
	app := InitializeScheduler()
	app.Run()
}
