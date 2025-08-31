package realtime

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Routes() chi.Router {
	realtimeRouter := chi.NewRouter()
	realtimeRouter.Get("/ws", handlerCreateWSConnection)
	return realtimeRouter
}

func handlerCreateWSConnection(w http.ResponseWriter, r *http.Request) {
	slog.Info("websocket connection establishment", "step", "begin")

	// TODO: Fix this to check origin correctly, also install and use cors package
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("websocket connection establishment", "step", "upgrade", "error", err)
		return
	}

	err = ws.WriteMessage(1, []byte("Hello from server"))
	if err != nil {
		slog.Error("websocket connection establishment", "step", "hello", "error", err)
		return
	}

	slog.Info("websocket connection establishment", "step", "finished")
}
