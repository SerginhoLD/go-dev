package io

import (
	"context"
	"exampleapp/domain/eventdispatcher"
	"exampleapp/io/controller"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

type App struct {
	eventDispatcher      eventdispatcher.EventDispatcher
	homeController       *controller.HomeController
	getProductController *controller.GetProductController
}

func NewApp(
	eventDispatcher eventdispatcher.EventDispatcher,
	homeController *controller.HomeController,
	getProductController *controller.GetProductController,
) *App {
	return &App{
		eventDispatcher,
		homeController,
		getProductController,
	}
}

func (app *App) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", controller.NotFoundHandler)
	mux.HandleFunc("GET /{$}", app.homeController.ServeHTTP) // https://pkg.go.dev/net/http#hdr-Patterns-ServeMux
	mux.HandleFunc("GET /product/{id}", app.getProductController.ServeHTTP)
	mux.Handle("GET /metrics", promhttp.Handler())

	http.ListenAndServe(":8080", app.httpLogMiddleware(mux))
}

func (app *App) httpLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestId, _ := uuid.NewV7()
		w.Header().Set("X-Request-ID", requestId.String())
		logResponseWriter := &LogResponseWriter{http.StatusInternalServerError, w}
		logRequest := r.WithContext(context.WithValue(r.Context(), "X-Request-ID", requestId.String()))

		defer func() {
			app.eventDispatcher.Dispatch(logRequest.Context(), &controller.ResponseEvent{logRequest, logResponseWriter.StatusCode, time.Since(start)})
		}()

		next.ServeHTTP(logResponseWriter, logRequest)
	})
}

type LogResponseWriter struct {
	StatusCode     int
	responseWriter http.ResponseWriter
}

func (w *LogResponseWriter) Header() http.Header {
	return w.responseWriter.Header()
}

func (w *LogResponseWriter) WriteHeader(statusCode int) {
	w.responseWriter.WriteHeader(statusCode)
	w.StatusCode = statusCode
}

func (w *LogResponseWriter) Write(b []byte) (int, error) {
	return w.responseWriter.Write(b)
}
