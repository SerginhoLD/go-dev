package io

import (
	"context"
	"exampleapp/infrastructure/logger"
	"exampleapp/io/controller"
	"fmt"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type App struct {
	logger                  *slog.Logger
	metrics                 *logger.Metrics
	coverageController      *controller.CoverageController
	homeController          *controller.HomeController
	getProductController    *controller.GetProductController
	createProductController *controller.CreateProductController
}

func NewApp(
	logger *slog.Logger,
	metrics *logger.Metrics,
	coverageController *controller.CoverageController,
	homeController *controller.HomeController,
	getProductController *controller.GetProductController,
	createProductController *controller.CreateProductController,
) *App {
	return &App{
		logger,
		metrics,
		coverageController,
		homeController,
		getProductController,
		createProductController,
	}
}

func (app *App) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", controller.NotFoundHandler)
	mux.HandleFunc("GET /{$}", app.homeController.ServeHTTP) // https://pkg.go.dev/net/http#hdr-Patterns-ServeMux
	mux.HandleFunc("GET /products/{id}", app.getProductController.ServeHTTP)
	mux.HandleFunc("POST /products", app.createProductController.ServeHTTP)
	mux.Handle("GET /metrics", promhttp.Handler())
	mux.HandleFunc("GET /coverage", app.coverageController.ServeHTTP)

	http.ListenAndServe(":8080", app.httpLogMiddleware(mux))
}

func (app *App) httpLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logResponseWriter := &LogResponseWriter{http.StatusInternalServerError, w}
		requestId, _ := uuid.NewV7()
		w.Header().Set("X-Request-ID", requestId.String())
		logRequest := r.WithContext(context.WithValue(r.Context(), "X-Request-ID", requestId.String()))

		defer func() {
			app.logger.DebugContext(logRequest.Context(), fmt.Sprintf("%s", logRequest.Pattern), slog.Int("statusCode", logResponseWriter.StatusCode))
			app.metrics.Add("app_http_request_duration_ms", float64(time.Since(start).Milliseconds()), logRequest.Pattern, strconv.Itoa(logResponseWriter.StatusCode))
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
