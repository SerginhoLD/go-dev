package main

import (
	"net/http"
)

func main() {
	app := InitializeApp()

	http.HandleFunc("/hello", app.homeController.ServeHTTP)
	http.HandleFunc("/headers", app.headersController.ServeHTTP)

	http.ListenAndServe(":8080", nil)
}
