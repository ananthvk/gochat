package realtime

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const maxEventsHub = 100

// TODO: For now Hub contains all the clients, rooms have not yet been implemented

// Hub manages a set of websocket connections
// It handles passing of events between clients and the rest of the system
type Hub struct {
	Clients map[uuid.UUID]*Client
	Events  chan []byte
}

func NewHub() *Hub {
	return &Hub{
		Clients: make(map[uuid.UUID]*Client),
		Events:  make(chan []byte, maxEventsHub),
	}
}

// RunEventLoop starts the event loop of the Hub. Note: This function must be called in a separate goroutine.
// It waits on events channel, and processes the events it received from the clients.
func (h *Hub) RunEventLoop() {
	slog.Info("started hub event loop")
	for {
		event := <-h.Events
		slog.Info("processed event", "size", len(event), "event", string(event))
	}
}

// AddConnection adds the websocket connection to the hub and returns the client ID.
// It also starts the Reader and Writer loop as two goroutines for the connection
func (h *Hub) AddConnection(connection *websocket.Conn) uuid.UUID {
	client := NewClient(connection, h)
	h.Clients[client.ID] = client
	go client.ReaderLoop()
	return client.ID
}
