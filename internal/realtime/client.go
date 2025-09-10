package realtime

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Maximum number of events that can remain unsent
const maxClientOutgoingSize = 100

const (
	// Maximum time allowed to write message to client
	maxWriteWait = 15 * time.Second

	// Duration of time the server waits for a PONG reply
	pongWait = 60 * time.Second

	// pingInterval defines how often ping messages are sent to clients
	// It has to be sent before the pong timeout
	pingInterval = (pongWait * 9) / 10

	// Maximum size of a message in bytes
	maxMessageSize = 4096
)

// client represents a websocket connection to the server.
// It stores additional fields such as ConnectedAt,
// and an UUID to uniquely identify this connection
// Outgoing is a buffered channel, and is used by the Hub to send Events to client
type client struct {
	ID          uuid.UUID
	Connection  *websocket.Conn
	ConnectedAt time.Time
	Outgoing    chan wsDataPacket
	Hub         *hub
}

// newClient creates a client from the passed websocket connection.
// It sets an unique ID to this connection, and creates the outgoing channel
func newClient(connection *websocket.Conn, hub *hub) *client {
	return &client{
		ID:          uuid.New(),
		Connection:  connection,
		ConnectedAt: time.Now().UTC(),
		Outgoing:    make(chan wsDataPacket, maxClientOutgoingSize),
		Hub:         hub,
	}
}

// ReaderLoop must be run in a separate goroutine. This function runs until the connection is terminated.
// It reads events from the client, and passes those events to the Hub for further processing
func (c *client) ReaderLoop() {
	defer func() {
		c.Hub.control <- hubConnectionUnregistered{Client: c}
		err := c.Connection.Close()
		if err != nil {
			slog.Error("error while closing websocket", "error", err)
		}
		slog.Info("closed websocket", "clientId", c.ID)
	}()

	c.Connection.SetReadLimit(maxMessageSize)
	c.Connection.SetReadDeadline(time.Now().Add(pongWait))
	// On receivng a pong message, extend the read deadline
	// This helps remove dead clients, i.e. if a client does not respond to a ping sent within the wait time, the read times out
	c.Connection.SetPongHandler(func(string) error { c.Connection.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		messageType, p, err := c.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("websocket read failed", "clientId", c.ID, "error", err)
			}
			return
		}

		var packet wsDataPacket

		if err := json.Unmarshal(p, &packet); err != nil {
			slog.Warn("json unmarshalling of packet failed", "clientId", c.ID, "messageType", messageType, "size", len(p), "error", err)
			continue
		}

		// Send the packet to the hub for further processing
		// If the Hub Events is full, drop the packet, so that the client retransmits it again
		select {
		case c.Hub.events <- hubDataReceived{Client: c, Packet: packet}:
			slog.Info("message enqueued to hub", "clientId", c.ID, "messageType", messageType, "size", len(p))
		default:
			slog.Warn("hub events channel full, dropped packet", "clientId", c.ID, "messageType", messageType, "size", len(p))
		}
	}
}

// WriterLoop must be run in a separate goroutine. This function runs until the connection is terminated.
// It waits for outgoing messages (in the Outgoing channel) and sends them to the client
func (c *client) WriterLoop() {
	// TODO: Optimization: Group multiple messages into a single message
	ticker := time.NewTicker(pingInterval)

	defer func() {
		ticker.Stop()
		c.Connection.Close()
	}()

	for {
		select {
		case message, ok := <-c.Outgoing:
			c.Connection.SetWriteDeadline(time.Now().Add(maxWriteWait))
			if !ok {
				// The hub has removed the client, send a close message
				c.Connection.WriteMessage(websocket.CloseMessage, []byte{})
				slog.Info("sent close message", "clientId", c.ID)
				return
			}
			b, err := json.Marshal(message)
			if err != nil {
				slog.Error("could not marshal outgoing message to bytes", "clientId", c.ID, "size", len(b), "error", err)
				continue
			}

			err = c.Connection.WriteMessage(websocket.TextMessage, b)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					slog.Error("message delivery failed", "clientId", c.ID, "size", len(b), "error", err)
				}
				return
			}
			slog.Info("message delivery successful", "clientId", c.ID, "size", len(b))
		case <-ticker.C:
			c.Connection.SetWriteDeadline(time.Now().Add(maxWriteWait))
			err := c.Connection.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				slog.Info("ping message to client failed", "clientId", c.ID)
				return
			}
		}
	}
}
