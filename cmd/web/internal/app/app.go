package app

import (
	"context"
	"database/sql"
	"errors"
	"exampleapp/internal/infrastructure/logger"
	"exampleapp/internal/infrastructure/postgres"
	"fmt"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"
)

type App struct {
	logger                  *slog.Logger
	metrics                 *logger.Metrics
	conn                    *postgres.Conn
	coverageController      *CoverageController
	homeController          *HomeController
	getProductController    *GetProductController
	createProductController *CreateProductController
}

func New(
	logger *slog.Logger,
	metrics *logger.Metrics,
	conn *postgres.Conn,
	coverageController *CoverageController,
	homeController *HomeController,
	getProductController *GetProductController,
	createProductController *CreateProductController,
) *App {
	return &App{
		logger,
		metrics,
		conn,
		coverageController,
		homeController,
		getProductController,
		createProductController,
	}
}

func (app *App) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", NotFoundHandler)
	mux.HandleFunc("GET /{$}", app.homeController.ServeHTTP) // https://pkg.go.dev/net/http#hdr-Patterns-ServeMux
	mux.HandleFunc("GET /products/{id}", app.getProductController.ServeHTTP)
	mux.Handle("POST /products", app.transactionMiddleware(http.HandlerFunc(app.createProductController.ServeHTTP)))
	mux.Handle("GET /metrics", promhttp.Handler())
	mux.HandleFunc("GET /coverage", app.coverageController.ServeHTTP)

	app.logger.Debug(fmt.Sprintf("server: start (env=%s)", os.Getenv("APP_ENV")))

	if err := http.ListenAndServe(":8080", app.logMiddleware(mux)); err != nil {
		app.logger.Error(fmt.Sprintf("server: %s", err.Error()))
	}
}

func (app *App) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logResponseWriter := &LogResponseWriter{http.StatusInternalServerError, w}
		requestId, _ := uuid.NewV7()
		w.Header().Set("X-Request-ID", requestId.String())
		*r = *r.WithContext(context.WithValue(r.Context(), "X-Request-ID", requestId.String()))

		// изначально r.Pattern не заполнен
		defer func() {
			app.logger.DebugContext(r.Context(), fmt.Sprintf(`server: response "%s"`, r.Pattern), slog.Int("statusCode", logResponseWriter.StatusCode))
			app.metrics.Add("app_http_request_duration_ms", float64(time.Since(start).Milliseconds()), r.Pattern, strconv.Itoa(logResponseWriter.StatusCode))
		}()

		next.ServeHTTP(logResponseWriter, r)
	})
}

func (app *App) transactionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		withTx := false
		switch r.Method {
		case "POST", "PATCH", "DELETE":
			withTx = true
		}

		if !withTx {
			next.ServeHTTP(w, r)
			return
		}

		tx, err := app.conn.DB().BeginTx(r.Context(), nil)

		if err != nil {
			HttpJsonError(w, err.Error(), http.StatusInternalServerError)
			return
		}

		app.logger.DebugContext(r.Context(), fmt.Sprintf(`sql: begin "%s %s"`, r.Method, r.RequestURI))
		*r = *r.WithContext(context.WithValue(r.Context(), "*sql.Tx", tx))

		defer func() {
			errRollback := tx.Rollback()

			if errRollback == nil {
				app.logger.DebugContext(r.Context(), fmt.Sprintf(`sql: rollback "%s %s"`, r.Method, r.RequestURI))
			} else if !errors.Is(errRollback, sql.ErrTxDone) {
				panic(errRollback)
			}
		}()

		defer func() {
			if rec := recover(); rec != nil {
				w.(*LogResponseWriter).StatusCode = http.StatusInternalServerError
				panic(rec) // вызываются все предыдущее defer
			}

			if w.(*LogResponseWriter).StatusCode < 400 {
				errCommit := tx.Commit()

				if errCommit != nil {
					panic(errCommit)
				}

				app.logger.DebugContext(r.Context(), fmt.Sprintf(`sql: commit "%s %s"`, r.Method, r.RequestURI))
			}
		}()

		next.ServeHTTP(w, r)
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
