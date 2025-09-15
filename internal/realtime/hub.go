package realtime

import (
	"encoding/json"
	"fmt"
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
	rooms   map[uuid.UUID]*hubRoom
	events  chan event
	control chan event
}

func newHub() *hub {
	return &hub{
		clients: make(map[uuid.UUID]*client),
		rooms:   make(map[uuid.UUID]*hubRoom),
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
		if payload.RoomID == uuid.Nil {
			slog.Error("roomID field missing", "from", e.ClientId, "message_type", "chat_message")
			return
		}
		h.broadcastMessageExceptSender(e.ClientId, payload)
	default:
		slog.Warn("invalid message type", "type", message.Type)
	}
}

// broadcastMessageExceptSender broadcasts a chat message to all connected clients in the same room
// except the client who sent it.
// If the outgoing channel of a client is full, the message is dropped and a warning is logged.
func (h *hub) broadcastMessageExceptSender(id uuid.UUID, payload wsChatMessage) {
	sender, ok := h.clients[id]
	if !ok {
		slog.Warn("broadcast failed since client does not exist", "id", id)
		return
	}
	room := h.rooms[payload.RoomID]
	if room == nil {
		slog.Warn("room does not exist, dropping message", "from", sender.ID, "room", payload.RoomID)
		return
	}

	for clientId := range room.Clients {
		client := h.clients[clientId]
		if client == nil {
			slog.Info("client does not exist anymore")
			// Also remove the client from the room
			delete(room.Clients, clientId)
			continue
		}

		if client.ID == sender.ID {
			continue
		}

		wsMessage := wsDataPacket{
			Type: "chat_message",
			Payload: toRaw(wsChatMessage{
				RoomID:  payload.RoomID,
				Message: payload.Message,
			}),
		}
		select {
		case client.Outgoing <- wsMessage:
			slog.Info("sent data", "from", sender.ID, "to", client.ID, "size", len(payload.Message), "room", payload.RoomID)
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
	case hubRoomCreated:
		h.processRoomCreateEvent(e)
	case hubRoomJoined:
		h.processRoomJoinEvent(e)
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

	// Send the welcome/connected message to the client with the client's id
	// The client requires this id to perform HTTP requests with the server
	h.clients[e.ClientId].Outgoing <- wsDataPacket{
		Type: "welcome",
		Payload: json.RawMessage(
			fmt.Appendf(nil, `{"id":"%s"}`, e.ClientId.String()),
		),
	}

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

		// Remove the client from all the rooms they are connected to
		for roomId := range client.Rooms {
			room := h.rooms[roomId]
			if room == nil {
				continue
			}
			slog.Info("removed client from room", "clientId", e.ClientId, "roomId", room.Id)
			delete(room.Clients, e.ClientId)
		}
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

func (h *hub) processRoomJoinEvent(e hubRoomJoined) {
	room, ok := h.rooms[e.roomId]
	if !ok {
		slog.Warn("join room(hub) failed", "reason", "room does not exist", "clientId", e.ClientId, "roomId", e.roomId)
		return
	}

	_, ok = room.Clients[e.ClientId]
	if ok {
		slog.Warn("join room(hub) failed", "reason", "client already part of room", "clientId", e.ClientId, "roomId", e.roomId)
		return
	}

	client := h.clients[e.ClientId]
	if client == nil {
		slog.Warn("join room(hub) failed", "reason", "client does not exist", "clientId", e.ClientId, "roomId", e.roomId)
		return
	}

	room.Clients[e.ClientId] = struct{}{}
	client.Rooms[e.roomId] = struct{}{}
	slog.Info("join room(hub) successful", "clientId", client.ID, "roomId", e.roomId)
}

func (h *hub) processRoomCreateEvent(e hubRoomCreated) {
	_, ok := h.rooms[e.roomId]
	if ok {
		slog.Info("not creating new room(hub) since it already exists", "id", e.roomId)
		return
	}
	h.rooms[e.roomId] = newHubRoom(e.roomId)
	slog.Info("created new room(hub)", "id", e.roomId)
}
