package realtime

import (
	"encoding/json"
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
			h.processEvent(event)
		case event := <-h.control:
			h.processControlEvent(event)
		}
	}
}

// processEvent processes a normal event (i.e. one that is not a control event).
// For now only hubDataReceived event is processed by this function
func (h *hub) processEvent(event event) {
	switch e := event.(type) {
	case hubDataReceived:
		h.processDataReceivedEvent(e)
	default:
		slog.Error("internal error", "reason", "unknown event")
		panic("unknown event")
	}
}

// processDataReceivedEvent handles incoming data packets from websocket clients.
// It unmarshals the packet payload based on the message type and performs the appropriate action.
// For "chat_message" type, it deserializes the payload into a wsChatMessage struct and broadcasts
// the message to all connected clients except the sender. Invalid message types or malformed
// JSON payloads are logged as warnings.
func (h *hub) processDataReceivedEvent(e hubDataReceived) {
	message := e.Packet
	switch message.Type {
	case "chat_message":
		var payload wsChatMessage
		if err := json.Unmarshal(message.Payload, &payload); err != nil {
			slog.Warn("malformed json payload", "from", e.ClientId, "message_type", "chat_message", "error", err)
			return
		}
		h.broadcastMessageExceptSender(e.ClientId, payload)
	default:
		slog.Warn("invalid message type", "type", message.Type)
	}
}

// broadcastMessageExceptSender broadcasts a chat message to all connected clients except the client who sent it.
// If the outgoing channel of a client is full, the message is dropped and a warning is logged.
func (h *hub) broadcastMessageExceptSender(id uuid.UUID, payload wsChatMessage) {
	sender, ok := h.clients[id]
	if !ok {
		slog.Warn("broadcast failed since client does not exist", "id", id)
	}
	for _, client := range h.clients {
		if client.ID == sender.ID {
			continue
		}
		wsMessage := wsDataPacket{
			Type: "chat_message",
			Payload: toRaw(wsChatMessage{
				Message: payload.Message,
			}),
		}
		select {
		case client.Outgoing <- wsMessage:
			slog.Info("sent data", "from", sender.ID, "to", client.ID, "size", len(payload.Message))
		default:
			slog.Warn("client outgoing channel full, dropping message", "from", sender.ID, "to", client.ID)
		}
	}
}

func toRaw(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}

// processControlEvent handles control events for the hub, processing connection
// registration and unregistration events. It routes the event to the appropriate
// handler based on the event type. If an unknown event type is received, it logs
// an error and panics to indicate an internal programming error.
func (h *hub) processControlEvent(event event) {
	switch e := event.(type) {
	case hubConnectionRegistered:
		h.processRegisterEvent(e)
	case hubConnectionUnregistered:
		h.processUnregisterEvent(e)
	default:
		slog.Error("internal error", "reason", "unknown control event")
		panic("unknown control event")
	}
}

// processRegisterEvent handles a register event. This event is generated when a new connection is created.
// It also starts a reader and writer loop for the new connection, in two new goroutines.
func (h *hub) processRegisterEvent(e hubConnectionRegistered) {
	client := newClient(e.ClientId, e.Connection, h)
	h.clients[e.ClientId] = client
	go client.ReaderLoop()
	go client.WriterLoop()
	slog.Info("processed register event", "clientId", e.ClientId)
}

// processUnregisterEvent handles an unregister event. This event is generated when a connection is closed
// when the client disconnects. It deletes the client from the clients map.
func (h *hub) processUnregisterEvent(e hubConnectionUnregistered) {
	// Check if the client is active
	if client, ok := h.clients[e.ClientId]; ok {
		delete(h.clients, e.ClientId)
		close(client.Outgoing)
		slog.Info("processed unregister event", "clientId", e.ClientId)
	} else {
		slog.Warn("unregister failed", "clientId", e.ClientId, "reason", "client with specified id does not exist")
	}
}

// addConnection adds the websocket connection to the hub and returns the client ID.
// It also starts the Reader and Writer loop as two goroutines for the connection
func (h *hub) addConnection(connection *websocket.Conn) uuid.UUID {
	id := uuid.New()
	h.control <- hubConnectionRegistered{ClientId: id, Connection: connection}
	return id
}
