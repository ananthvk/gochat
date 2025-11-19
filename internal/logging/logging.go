package logging

import (
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/ananthvk/gochat/internal/config"
	"github.com/go-chi/httplog/v3"
	"github.com/go-chi/traceid"
	"github.com/golang-cz/devslog"
)

func CreateLoggerAndRequestLoggerMiddleware(cfg *config.Config) (*slog.Logger, func(http.Handler) http.Handler) {
	isLocalhost := cfg.Env == "development"
	logFormat := httplog.SchemaECS.Concise(isLocalhost)
	logger := slog.New(logHandler(isLocalhost, &slog.HandlerOptions{
		AddSource:   !isLocalhost,
		ReplaceAttr: logFormat.ReplaceAttr,
	}))
	middleware := createRequestLoggerMiddleware(logger, logFormat, cfg)
	return logger, middleware
}

func logHandler(isLocalhost bool, handlerOpts *slog.HandlerOptions) slog.Handler {
	if isLocalhost {
		// Pretty logs for development.
		return devslog.NewHandler(os.Stdout, &devslog.Options{
			SortKeys:           true,
			MaxErrorStackTrace: 5,
			MaxSlicePrintSize:  20,
			HandlerOptions:     handlerOpts,
		})
	}

	// JSON logs for production with "traceId".
	return traceid.LogHandler(
		slog.NewJSONHandler(os.Stdout, handlerOpts),
	)
}

// Returns a middleware that logs requests
func createRequestLoggerMiddleware(logger *slog.Logger, logFormat *httplog.Schema, cfg *config.Config) func(http.Handler) http.Handler {
	isLocalhost := cfg.Env == "development"
	return httplog.RequestLogger(logger, &httplog.Options{
		// Level defines the verbosity of the request logs:
		// slog.LevelDebug - log all responses (incl. OPTIONS)
		// slog.LevelInfo  - log all responses (excl. OPTIONS)
		// slog.LevelWarn  - log 4xx and 5xx responses only (except for 429)
		// slog.LevelError - log 5xx responses only
		Level: slog.LevelInfo,

		// Log attributes using given schema/format.
		Schema: logFormat,

		// RecoverPanics recovers from panics occurring in the underlying HTTP handlers
		// and middlewares. It returns HTTP 500 unless response status was already set.
		//
		// NOTE: Panics are logged as errors automatically, regardless of this setting.
		RecoverPanics: true,

		// Filter out some request logs.
		Skip: func(req *http.Request, respStatus int) bool {
			return respStatus == 404 || respStatus == 405
		},

		// Select request/response headers to be logged explicitly.
		LogRequestHeaders:  []string{"Origin"},
		LogResponseHeaders: []string{},

		// You can log request/request body conditionally. Useful for debugging.
		LogRequestBody:  isDebugHeaderSet,
		LogResponseBody: isDebugHeaderSet,

		// Log all requests with invalid payload as curl command.
		LogExtraAttrs: func(req *http.Request, reqBody string, respStatus int) []slog.Attr {
			// Only print curl command in development mode
			if !isLocalhost {
				return nil
			}
			if respStatus == 400 || respStatus == 422 {
				req.Header.Del("Authorization")
				// If it's signup URL, don't log curl command since it contains the password in the body
				if strings.Contains(req.URL.String(), "signup") {
					return nil
				}
				return []slog.Attr{slog.String("curl", httplog.CURL(req, reqBody))}
			}
			return nil
		},
	})
}

func isDebugHeaderSet(r *http.Request) bool {
	return r.Header.Get("Debug") == "TRUE"
}
