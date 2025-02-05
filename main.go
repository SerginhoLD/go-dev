package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/hello", InitializeHomeController().ServeHTTP)
	http.HandleFunc("/headers", InitializeHeadersController().ServeHTTP)

	http.ListenAndServe(":8080", nil)
}
