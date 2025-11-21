package main

import (
	"context"
	"fmt"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ananthvk/gochat/internal"
	"github.com/ananthvk/gochat/internal/app"
	"github.com/ananthvk/gochat/internal/auth"
	"github.com/ananthvk/gochat/internal/config"
	"github.com/ananthvk/gochat/internal/logging"
	mid "github.com/ananthvk/gochat/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/traceid"
)

const appVersion = "0.0.1"

func main() {
	startTime := time.Now().UTC()
	// Load configuration from environment and dotfiles
	config.LoadEnv()
	cfg, err := config.ParseConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s: %s\n", "Invalid environment variable", err.Error())
		os.Exit(1)
	}

	// Set default logger to httplog
	logger, requestLoggerMiddleware := logging.CreateLoggerAndRequestLoggerMiddleware(cfg)
	if cfg.Env != "development" {
		logger = logger.With(
			slog.String("app", "gochat"),
			// TODO: Get the version number and running environment automatically
			slog.String("version", appVersion),
			slog.String("env", cfg.Env),
		)
	}
	slog.SetDefault(logger)
	slog.SetLogLoggerLevel(slog.LevelError)

	if cfg.Env == "development" {
		// If in development environment, also log all config values
		slog.Info("current configuration", "cfg", cfg)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(traceid.Middleware)
	router.Use(middleware.Recoverer)
	router.Use(requestLoggerMiddleware)
	router.Use(middleware.Heartbeat("/ping"))

	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".css", "text/css")

	slog.Info("finished loading middlewares")

	ctx := context.Background()

	app, err := app.NewApp(ctx, cfg, appVersion)
	if err != nil {
		slog.Error("exiting since app could not be created", "error", err)
	}
	app.StartTime = startTime
	app.RealtimeService.StartHubEventLoop()
	defer app.DatabaseService.Pool.Close()

	middleware := mid.Middlewares{
		Authenticate: auth.AuthMiddleware(app.TokenService),
	}

	router.Mount("/api/v1/", internal.Routes(app, middleware))

	fs := http.FileServer(http.Dir("./static"))
	router.Handle("/*", fs)

	server := &http.Server{
		Addr:    cfg.Host + ":" + strconv.Itoa(cfg.Port),
		Handler: router,
	}

	slog.Info("server listening", "address", server.Addr)
	slog.Error("server quit", "error", server.ListenAndServe())
}
