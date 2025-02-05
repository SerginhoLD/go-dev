package main

import (
	"example.com/m/io/controller"
	"fmt"
	"net/http"
	"os"
	//"log"
	"log/slog"
)

/*
func hello(w http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(w, "hello 2\n")

	//log.Println("standard logger")

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	logger.Info("hello again", "key", "val", "age", 25)
}*/

func headers(w http.ResponseWriter, req *http.Request) {

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	homeController := controller.NewHomeController(*logger)

	http.HandleFunc("/hello", homeController.Index)
	http.HandleFunc("/headers", headers)

	http.ListenAndServe(":8080", nil)
}
