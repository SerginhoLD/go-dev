package main

import (
	"exampleapp/cmd/consumer/internal"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env.local")
	godotenv.Load(".env")
	app := internal.InitializeApp(internal.GroupName(os.Args[1]), internal.StreamName(os.Args[2]))
	app.Run()
}
