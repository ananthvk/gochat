package internal

import (
	"github.com/ananthvk/gochat/internal/realtime"
	"github.com/go-chi/chi/v5"
)

func Routes() chi.Router {
	router := chi.NewRouter()
	router.Mount("/realtime", realtime.Routes())
	return router
}
