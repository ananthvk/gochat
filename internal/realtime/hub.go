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
	Events  chan Event
	Control chan Event
}

func NewHub() *Hub {
	return &Hub{
		Clients: make(map[uuid.UUID]*Client),
		Events:  make(chan Event, maxEventsHub),
		Control: make(chan Event),
	}
}

// RunEventLoop starts the event loop of the Hub. Note: This function must be called in a separate goroutine.
// It waits on events channel, and processes the events it received from the clients.
// TODO: NOTE: Since both kind of events (control, and data events) are multiplexed and processed in the same goroutine,
// it may lead to starvation. Research/Identify some method to prevent starvation.
func (h *Hub) RunEventLoop() {
	slog.Info("started hub event loop")
	for {
		select {
		case event := <-h.Events:
			switch e := event.(type) {
			case DataEvent:
				slog.Info("processed data event", "size", len(e.Data), "payload", string(e.Data))
			default:
				slog.Error("internal error", "reason", "unknown event")
				panic("unknown event")
			}
		case event := <-h.Control:
			switch e := event.(type) {
			case RegisterConnectionEvent:
				h.Clients[e.Client.ID] = e.Client
				go e.Client.ReaderLoop()
				slog.Info("processed register event", "clientId", e.Client.ID)
			case UnregisterConnectionEvent:
				// Check if the client is active
				if _, ok := h.Clients[e.Client.ID]; ok {
					delete(h.Clients, e.Client.ID)
					close(e.Client.Outgoing)
					slog.Info("processed unregister event", "clientId", e.Client.ID)
				} else {
					slog.Warn("unregister failed", "clientId", e.Client.ID, "reason", "client with specified id does not exist")
				}
			default:
				slog.Error("internal error", "reason", "unknown control event")
				panic("unknown control event")
			}
		}
	}
}

// AddConnection adds the websocket connection to the hub and returns the client ID.
// It also starts the Reader and Writer loop as two goroutines for the connection
func (h *Hub) AddConnection(connection *websocket.Conn) uuid.UUID {
	client := NewClient(connection, h)
	h.Control <- RegisterConnectionEvent{Client: client}
	return client.ID
}
