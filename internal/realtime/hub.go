package realtime

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const maxEventsHub = 100

// TODO: For now Hub contains all the clients, rooms have not yet been implemented

// hub manages a set of websocket connections
// It handles passing of events between clients and the rest of the system
type hub struct {
	clients map[uuid.UUID]*client
	events  chan event
	control chan event
}

func newHub() *hub {
	return &hub{
		clients: make(map[uuid.UUID]*client),
		events:  make(chan event, maxEventsHub),
		control: make(chan event),
	}
}

// RunEventLoop starts the event loop of the Hub. Note: This function must be called in a separate goroutine.
// It waits on events channel, and processes the events it received from the clients.
// TODO: NOTE: Since both kind of events (control, and data events) are multiplexed and processed in the same goroutine,
// it may lead to starvation. Research/Identify some method to prevent starvation.
func (h *hub) RunEventLoop() {
	slog.Info("started hub event loop")
	for {
		select {
		case event := <-h.events:
			switch e := event.(type) {
			case dataEvent:
				// For now, send the message to all connected clients
				for _, client := range h.clients {
					if client.ID == e.Client.ID {
						continue
					}
					select {
					case client.Outgoing <- e.Data:
						slog.Info("sent data", "from", e.Client.ID, "to", client.ID, "size", len(e.Data))
					default:
						slog.Warn("client outgoing channel full, dropping message", "from", e.Client.ID, "to", client.ID, "size", len(e.Data))
					}
				}
				slog.Info("processed data event", "size", len(e.Data), "payload", string(e.Data))
			default:
				slog.Error("internal error", "reason", "unknown event")
				panic("unknown event")
			}
		case event := <-h.control:
			switch e := event.(type) {
			case registerConnectionEvent:
				h.clients[e.Client.ID] = e.Client
				go e.Client.ReaderLoop()
				go e.Client.WriterLoop()
				slog.Info("processed register event", "clientId", e.Client.ID)
			case unregisterConnectionEvent:
				// Check if the client is active
				if _, ok := h.clients[e.Client.ID]; ok {
					delete(h.clients, e.Client.ID)
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

// addConnection adds the websocket connection to the hub and returns the client ID.
// It also starts the Reader and Writer loop as two goroutines for the connection
func (h *hub) addConnection(connection *websocket.Conn) uuid.UUID {
	client := newClient(connection, h)
	h.control <- registerConnectionEvent{Client: client}
	return client.ID
}
