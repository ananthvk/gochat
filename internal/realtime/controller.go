package realtime

import (
	"log/slog"
	"net/http"
	"strings"

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
	realtimeRouter.Post("/room", func(w http.ResponseWriter, r *http.Request) { handleCreateRoom(rt, w, r) })
	realtimeRouter.Get("/room/by-name/{name}", func(w http.ResponseWriter, r *http.Request) { handleGetRoomByName(rt, w, r) })
	realtimeRouter.Get("/room", func(w http.ResponseWriter, r *http.Request) { handleGetRooms(rt, w, r) })
	return realtimeRouter
}

// handleGetRooms returns all the rooms available on the server
func handleGetRooms(rt *RealtimeService, w http.ResponseWriter, _ *http.Request) {
	rooms := rt.ListRooms()
	helpers.RespondWithJSON(w, http.StatusOK, rooms)
}

// handleGetRoomByName returns the room associated with a room name. For now,
// rooms are uniquely identifiable by their room name.
func handleGetRoomByName(rt *RealtimeService, w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if strings.TrimSpace(name) == "" {
		helpers.RespondWithError(w, http.StatusBadRequest, "url parameter missing", "required URL parameter 'name' not provided")
		return
	}

	room := rt.GetRoomByName(name)
	if room == nil {
		helpers.RespondWithError(w, http.StatusNotFound, "room not found", "room with given name not found")
		return
	}
	helpers.RespondWithJSON(w, http.StatusOK, room)
}

// handleCreateRoom creates a new room and returns the ID of the room to the client
func handleCreateRoom(rt *RealtimeService, w http.ResponseWriter, r *http.Request) {
	createRoomRequest := struct {
		Name string `json:"name"`
	}{}
	err := helpers.ParseJSON(r.Body, &createRoomRequest, false)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, "post body malformed", err.Error())
		return
	}
	room := rt.CreateRoom(createRoomRequest.Name)
	helpers.RespondWithJSON(w, http.StatusCreated, room)
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
