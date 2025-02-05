package main

import (
	"fmt"
	"net/http"
)

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {

	homeController := InitializeHomeController()

	http.HandleFunc("/hello", homeController.Index)
	http.HandleFunc("/headers", headers)

	http.ListenAndServe(":8080", nil)
}
