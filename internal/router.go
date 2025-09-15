package internal

import (
	"github.com/ananthvk/gochat/internal/app"
	"github.com/ananthvk/gochat/internal/realtime"
	"github.com/go-chi/chi/v5"
)

func Routes(app *app.App) chi.Router {
	router := chi.NewRouter()
	router.Mount("/realtime", realtime.Routes(app.RealtimeService))
	return router
}
