package testutils

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/ananthvk/gochat/internal"
	"github.com/ananthvk/gochat/internal/app"
	"github.com/ananthvk/gochat/internal/realtime"
	"github.com/go-chi/chi/v5"
)

func NewTestServer(t *testing.T, ctx context.Context) (*app.App, *chi.Mux) {
	t.Helper()
	router := chi.NewRouter()
	app := &app.App{
		Ctx:             ctx,
		RealtimeService: realtime.NewRealtimeService(ctx),
	}
	app.RealtimeService.StartHubEventLoop()
	router.Mount("/api/v1/", internal.Routes(app))
	return app, router
}

func NewTestServerWithCancel(t *testing.T) (*app.App, *httptest.Server, context.CancelFunc) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	app, router := NewTestServer(t, ctx)
	srv := httptest.NewServer(router)
	return app, srv, cancel
}
