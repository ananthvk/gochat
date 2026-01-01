package realtime

import (
	"context"
	"log/slog"

	"github.com/gorilla/websocket"
	"github.com/oklog/ulid/v2"
)

const maxEventsHub = 100

type clientSet = map[ulid.ULID]struct{}

type event any

type broadcastEvent struct {
	targetRoom ulid.ULID
	payload    []byte
}

type registerClientEvent struct {
	conn     *websocket.Conn
	userId   ulid.ULID
	clientId ulid.ULID
}

type unregisterClientEvent struct {
	clientId ulid.ULID
}

type createRoomsAndAddClientEvent struct {
	clientId ulid.ULID
	roomIds  []ulid.ULID
}

// hub manages a set of websocket connections
// It handles routing of messages
type hub struct {
	clients map[ulid.ULID]*client
	// rooms map room id to a set of clients
	rooms   map[ulid.ULID]clientSet
	events  chan event
	control chan event
}

func newHub() *hub {
	return &hub{
		clients: make(map[ulid.ULID]*client),
		rooms:   make(map[ulid.ULID]clientSet),
		events:  make(chan event, maxEventsHub),
		control: make(chan event),
	}
}

// RunEventLoop starts the event loop of the Hub. Note: This function must be called in a separate goroutine.
// It waits on events channel, and processes the events it received from the clients.
// TODO: NOTE: Since both kind of events (control, and data events) are multiplexed and processed in the same goroutine,
// it may lead to starvation. Research/Identify some method to prevent starvation.
func (h *hub) RunEventLoop(ctx context.Context) {
	slog.Info("started hub event loop")
	for {
		select {
		case event := <-h.events:
			h.processEvent(event)
		case event := <-h.control:
			h.processControlEvent(event)
		case <-ctx.Done():
			slog.Info("unregistering all connected clients")
			for _, client := range h.clients {
				e := unregisterClientEvent{clientId: client.ID}
				h.processUnregisterEvent(e)
			}
			slog.Info("stopped hub context loop", "reason", ctx.Err())
			return
		}
	}
}

// processEvent processes a normal event (i.e. one that is not a control event).
// For now only hubDataReceived event is processed by this function
func (h *hub) processEvent(ev event) {
	switch e := ev.(type) {
	case broadcastEvent:
		h.handleBroadcast(e)
	default:
		slog.Error("internal error", "reason", "unknown event")
		panic("unknown event")
	}
}

// handleBroadcast handles broadcasting of a message to connected clients in the targetRoom
// It does no marshalling of data, and the bytes are sent as received. In case of any error,
// the message is silently dropped. Since these events are created by the application, they are assumed
// to be correct, hence no error checking is done.
// If the outgoing channel of a client is full, the message is dropped silently
func (h *hub) handleBroadcast(e broadcastEvent) {
	room, ok := h.rooms[e.targetRoom]
	if !ok {
		slog.Warn("broadcast failed", "reason", "room does not exist", "id", e.targetRoom)
		return
	}
	for clientId := range room {
		client := h.clients[clientId]
		if client == nil {
			slog.Info("client does not exist anymore")
			// Note: Lazy deletion
			// Remove the client from the room
			delete(room, clientId)
			continue
		}
		select {
		case client.Outgoing <- e.payload:
		default:
		}
	}
}

// processControlEvent handles control events for the hub, processing connection
// registration and unregistration events. It routes the event to the appropriate
// handler based on the event type. If an unknown event type is received, it logs
// an error and panics to indicate an internal programming error.
func (h *hub) processControlEvent(event event) {
	switch e := event.(type) {
	case registerClientEvent:
		h.processRegisterEvent(e)
	case unregisterClientEvent:
		h.processUnregisterEvent(e)
	case createRoomsAndAddClientEvent:
		h.processCreateRoomAndJoinEvent(e)
	default:
		slog.Error("internal error", "reason", "unknown control event")
		panic("unknown control event")
	}
}

// processRegisterEvent handles a register event. This event is generated when a new connection is created.
// It also starts a reader and writer loop for the new connection, and spawns two goroutines for them
func (h *hub) processRegisterEvent(e registerClientEvent) {
	client := newClient(e.conn, e.userId, e.clientId, h)
	h.clients[client.ID] = client
	go client.ReaderLoop()
	go client.WriterLoop()
	slog.Info("processed register event", "clientId", client.ID)
}

// processUnregisterEvent handles an unregister event. This event is generated when a connection is closed
// when the client disconnects.
func (h *hub) processUnregisterEvent(e unregisterClientEvent) {
	// Check if the client is active
	if client, ok := h.clients[e.clientId]; ok {
		delete(h.clients, e.clientId)
		close(client.Outgoing)
		// Note: We are not removing the client from all the maps, since they get lazily deleted when a broadcast message is sent
		// Note: This might be an issue if say a rogue client repeatedly connects / disconnects causing the map to get full (when no messages are sent)
		slog.Info("processed unregister event", "clientId", e.clientId)
	}
}

func (h *hub) processCreateRoomAndJoinEvent(e createRoomsAndAddClientEvent) {
	if _, ok := h.clients[e.clientId]; !ok {
		return
	}
	for _, roomId := range e.roomIds {
		room, ok := h.rooms[roomId]
		if !ok {
			room = make(clientSet)
			h.rooms[roomId] = room
		}
		room[e.clientId] = struct{}{}
	}
	slog.Info("processed createRoomsAndAddClientEvent", "rooms", e.roomIds, "client", e.clientId)
}
