package main

import (
	"exampleapp/internal/app/consumer"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env.local")
	godotenv.Load(".env")
	app := consumer.Initialize(consumer.GroupName(os.Args[1]), consumer.StreamName(os.Args[2]))
	app.Run()
}
