package internal

import (
	"context"
	"errors"
	"exampleapp/internal/infrastructure/di"
	"exampleapp/internal/infrastructure/logger"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type App struct {
	coverageController *CoverageController
	homeController     *HomeController
}

func NewApp(
	coverageController *CoverageController,
	homeController *HomeController,
) *App {
	return &App{
		coverageController,
		homeController,
	}
}

func (app *App) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", NotFoundHandler) // https://pkg.go.dev/net/http#hdr-Patterns-ServeMux
	mux.HandleFunc("GET /{$}", app.homeController.ServeHTTP)
	mux.Handle("GET /metrics", promhttp.Handler())
	mux.HandleFunc("GET /coverage", app.coverageController.ServeHTTP)

	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("web/assets"))))
	mux.HandleFunc("GET /assets/{$}", NotFoundHandler)

	slog.Debug(fmt.Sprintf("server: start (env=%s, ver=%s)", os.Getenv("APP_ENV"), di.Version))

	if err := http.ListenAndServe(":8080", app.logMiddleware(mux)); !errors.Is(err, http.ErrServerClosed) {
		slog.Error(fmt.Sprintf("server: %s", err.Error()))
		os.Exit(1)
	}
}

func (app *App) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logResponseWriter := &logResponseWriter{http.StatusInternalServerError, w}
		requestId, _ := uuid.NewV7()
		w.Header().Set("X-Request-ID", requestId.String())
		*r = *r.WithContext(context.WithValue(r.Context(), "X-Request-ID", requestId.String()))

		// изначально r.Pattern не заполнен
		defer func() {
			slog.DebugContext(r.Context(), fmt.Sprintf(`server: response "%s"`, r.Pattern), slog.Int("statusCode", logResponseWriter.statusCode))
			logger.AddMetric("app_http_request_duration_ms", float64(time.Since(start).Milliseconds()), r.Pattern, strconv.Itoa(logResponseWriter.statusCode))
		}()

		next.ServeHTTP(logResponseWriter, r)
	})
}

type logResponseWriter struct {
	statusCode     int
	responseWriter http.ResponseWriter
}

func (w *logResponseWriter) Header() http.Header {
	return w.responseWriter.Header()
}

func (w *logResponseWriter) WriteHeader(statusCode int) {
	w.responseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func (w *logResponseWriter) Write(b []byte) (int, error) {
	return w.responseWriter.Write(b)
}
