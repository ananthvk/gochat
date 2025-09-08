package main

import (
	"log/slog"
	"mime"
	"net/http"
	"os"

	"github.com/ananthvk/gochat/internal"
	"github.com/ananthvk/gochat/internal/logging"
	"github.com/ananthvk/gochat/internal/realtime"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/traceid"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = "127.0.0.1"
	}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8000"
	}

	// Set default logger to httplog
	logger, requestLoggerMiddleware := logging.CreateLoggerAndRequestLoggerMiddleware()
	if k := os.Getenv("ENV"); k != "localhost" {
		logger = logger.With(
			slog.String("app", "gochat"),
			// TODO: Get the version number and running environment automatically
			slog.String("version", "v0.0.1"),
			slog.String("env", "production"),
		)
	}
	slog.SetDefault(logger)
	slog.SetLogLoggerLevel(slog.LevelError)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(traceid.Middleware)
	router.Use(middleware.Recoverer)
	router.Use(requestLoggerMiddleware)

	router.Use(middleware.Heartbeat("/ping"))

	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".css", "text/css")

	router.Mount("/api/v1/", internal.Routes())

	fs := http.FileServer(http.Dir("./static"))
	router.Handle("/*", fs)

	server := &http.Server{
		Addr:    host + ":" + port,
		Handler: router,
	}

	go realtime.DefaultHub.RunEventLoop()

	slog.Info("server listening", "address", server.Addr)
	slog.Error("server quit", "error", server.ListenAndServe())
}
