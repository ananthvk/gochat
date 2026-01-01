package realtime

import (
	"context"

	"github.com/ananthvk/gochat/internal/database"
	"github.com/gorilla/websocket"
	"github.com/oklog/ulid/v2"
)

type RealtimeService struct {
	clientHub *hub
	Db        *database.DatabaseService
}

func NewRealtimeService(ctx context.Context, db *database.DatabaseService) *RealtimeService {
	hub := newHub()
	go hub.RunEventLoop(ctx)
	return &RealtimeService{
		clientHub: hub,
		Db:        db,
	}
}

// RegisterConnection registers a new websocket connection and returns a connection id
func (r *RealtimeService) RegisterConnection(conn *websocket.Conn, userId ulid.ULID) ulid.ULID {
	clientId := ulid.Make()
	r.clientHub.control <- registerClientEvent{conn: conn, userId: userId, clientId: clientId}
	return clientId
}

func (r *RealtimeService) UnregisterConnection(clientId ulid.ULID) {
	r.clientHub.control <- unregisterClientEvent{clientId: clientId}
}

// AddConnectionToRooms creates the rooms from the specified list, if it exists, it's not created again, then the client is added to all those rooms
func (r *RealtimeService) AddConnectionToRooms(roomIds []ulid.ULID, clientId ulid.ULID) {
	r.clientHub.control <- createRoomsAndAddClientEvent{clientId: clientId, roomIds: roomIds}
}

func (r *RealtimeService) RemoveConnectionFromRoom(connectionId ulid.ULID, roomId ulid.ULID) {
	// TOOD: Implement this
	// Only needed when a user is removed from a group
}

func (r *RealtimeService) Broadcast(roomId ulid.ULID, message []byte) {
	r.clientHub.events <- broadcastEvent{targetRoom: roomId, payload: message}
}

// Other methods that are necesssary - A method to remove all connections associated with a client (incase of logout)
