package realtime

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Routes() chi.Router {
	realtimeRouter := chi.NewRouter()
	realtimeRouter.Get("/ws", handlerCreateWSConnection)
	return realtimeRouter
}

func handlerCreateWSConnection(w http.ResponseWriter, r *http.Request) {
	log.Printf("Creating websocket connection")
	fmt.Fprintf(w, "Hello %s", "world")
}
