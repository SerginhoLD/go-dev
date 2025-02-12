//go:build wireinject
// +build wireinject

package main

import (
	"exampleapp/domain/eventdispatcher"
	"exampleapp/domain/repository"
	"exampleapp/domain/usecase"
	eventdispatcherimpl "exampleapp/infrastructure/eventdispatcher"
	"exampleapp/infrastructure/logger"
	"exampleapp/infrastructure/postgres"
	repositoryimpl "exampleapp/infrastructure/repository"
	"exampleapp/io/controller"
	"fmt"
	"github.com/google/wire"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
)

type App struct {
	logger               *slog.Logger
	homeController       *controller.HomeController
	getProductController *controller.GetProductController
}

func NewApp(
	logger *slog.Logger,
	homeController *controller.HomeController,
	getProductController *controller.GetProductController,
) *App {
	return &App{
		logger,
		homeController,
		getProductController,
	}
}

func (app *App) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", app.homeController.ServeHTTP) // https://pkg.go.dev/net/http#hdr-Patterns-ServeMux
	mux.HandleFunc("GET /product/{id}", app.getProductController.ServeHTTP)

	mux.Handle("GET /metrics", promhttp.Handler())

	handler := app.httpLogMiddleware(mux)
	http.ListenAndServe(":8080", handler)
}

func (app *App) httpLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logResponseWriter := &LogResponseWriter{http.StatusOK, w}
		next.ServeHTTP(logResponseWriter, r)
		app.logger.Info(fmt.Sprintf("%s %s", r.Method, r.Pattern), slog.Int("statusCode", logResponseWriter.StatusCode))
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

func InitializeApp() *App {
	wire.Build(
		NewApp,
		logger.NewLogger,
		logger.NewLogListener,
		logger.NewMetricListener,
		postgres.NewConn,
		eventdispatcherimpl.New,
		wire.Bind(new(eventdispatcher.EventDispatcher), new(*eventdispatcherimpl.EventDispatcherImpl)),
		repositoryimpl.NewProductRepositoryImpl,
		wire.Bind(new(repository.ProductRepository), new(*repositoryimpl.ProductRepositoryImpl)),
		usecase.NewPaginateProductsUseCase,
		controller.NewHomeController,
		usecase.NewGetProductUseCase,
		controller.NewGetProductController,
	)

	return &App{}
}
