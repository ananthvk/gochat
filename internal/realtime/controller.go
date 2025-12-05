package realtime

import (
	"log/slog"
	"net/http"

	"github.com/ananthvk/gochat/internal/helpers"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Routes(rt *RealtimeService) chi.Router {
	realtimeRouter := chi.NewRouter()
	realtimeRouter.Get("/ws", func(w http.ResponseWriter, r *http.Request) { handlerCreateWSConnection(rt, w, r) })
	return realtimeRouter
}

func handlerCreateWSConnection(rt *RealtimeService, w http.ResponseWriter, r *http.Request) {

	// TODO: Fix this to check origin correctly, also install and use cors package
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		slog.ErrorContext(r.Context(), "websocket upgrade failed", "error", err)
		helpers.RespondWithError(w, http.StatusUpgradeRequired, "websocket upgrade failed", err.Error())
		return
	}

	client_id := rt.hub.addConnection(conn)
	slog.InfoContext(r.Context(), "websocket connection established", "address", conn.RemoteAddr(), "clientId", client_id)
}
