package internal

import (
	"net/http"

	"github.com/ananthvk/gochat/internal/app"
	"github.com/ananthvk/gochat/internal/group"

	"github.com/ananthvk/gochat/internal/health"
	"github.com/ananthvk/gochat/internal/realtime"
	"github.com/go-chi/chi/v5"
)

func Routes(app *app.App) chi.Router {
	router := chi.NewRouter()
	router.Mount("/realtime", realtime.Routes(app.RealtimeService))
	router.Mount("/group", group.Routes(app.GroupService, app.MessageService))
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) { health.HealthCheckHandler(app, w, r) })
	return router
}
