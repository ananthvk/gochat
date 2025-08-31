package logging

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/httplog/v3"
	"github.com/go-chi/traceid"
	"github.com/golang-cz/devslog"
)

func CreateLoggerAndRequestLoggerMiddleware() (*slog.Logger, func(http.Handler) http.Handler) {
	isLocalhost := os.Getenv("ENV") == "localhost"
	logFormat := httplog.SchemaECS.Concise(isLocalhost)
	logger := slog.New(logHandler(isLocalhost, &slog.HandlerOptions{
		AddSource:   !isLocalhost,
		ReplaceAttr: logFormat.ReplaceAttr,
	}))
	middleware := createRequestLoggerMiddleware(logger, logFormat)
	return logger, middleware
}

func logHandler(isLocalhost bool, handlerOpts *slog.HandlerOptions) slog.Handler {
	if isLocalhost {
		// Pretty logs for localhost development.
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
func createRequestLoggerMiddleware(logger *slog.Logger, logFormat *httplog.Schema) func(http.Handler) http.Handler {
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
			if respStatus == 400 || respStatus == 422 {
				req.Header.Del("Authorization")
				return []slog.Attr{slog.String("curl", httplog.CURL(req, reqBody))}
			}
			return nil
		},
	})
}

func isDebugHeaderSet(r *http.Request) bool {
	return r.Header.Get("Debug") == "TRUE"
}
