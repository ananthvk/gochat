package realtime

import (
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
	"github.com/oklog/ulid/v2"
)

const (
	// Maximum number of messages that can remain unsent for a client
	maxClientOutgoing = 100

	// Maximum time allowed to write mesage to client
	maxWriteWait = 15 * time.Second

	// The duration the server waits for a PONG reply
	pongWait = 60 * time.Second

	// pingInterval defines how often ping messages are sent to clients
	// It has to be sent before the pong timeout
	pingInterval = (pongWait * 9) / 10

	// Maximum size of a message
	maxMessageSize = 4096
)

type client struct {
	ID          ulid.ULID
	UserId      ulid.ULID
	Connection  *websocket.Conn
	ConnectedAt time.Time
	Outgoing    chan []byte
	clientHub   *hub
}

func newClient(conn *websocket.Conn, userId, clientId ulid.ULID, h *hub) *client {
	return &client{
		ID:          clientId,
		UserId:      userId,
		Connection:  conn,
		ConnectedAt: time.Now(),
		Outgoing:    make(chan []byte, maxClientOutgoing),
		clientHub:   h,
	}
}

// ReaderLoop must be run in a separate goroutine. This function runs until the connection is terminated.
// Since websockets are currently used only for notifications, the reader loop only listens for pong responses
// to keep the connection alive
func (c *client) ReaderLoop() {
	defer func() {
		c.clientHub.control <- unregisterClientEvent{clientId: c.ID}
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
		// Sliently drop messages
		// TODO: Later use this for typing / online indicators
		_, _, err := c.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("websocket read failed", "clientId", c.ID, "error", err, "connectionId", c.ID)
			}
			return
		}
	}
}

// WriterLoop must be run in a separate goroutine. This function runs until the connection is terminated.
// It waits for outgoing messages (in the Outgoing channel) and sends them to the client
func (c *client) WriterLoop() {
	// TODO: Optimization: Group multiple messages into a single message
	ticker := time.NewTicker(pingInterval)

	defer func() {
		c.clientHub.control <- unregisterClientEvent{clientId: c.ID}
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
			if err := c.Connection.WriteMessage(websocket.TextMessage, message); err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					slog.Error("message delivery failed", "clientId", c.ID, "size", len(message), "error", err)
				}
				return
			}
			slog.Info("message delivery successful", "clientId", c.ID, "size", len(message))
		case <-ticker.C:
			c.Connection.SetWriteDeadline(time.Now().Add(maxWriteWait))
			err := c.Connection.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				slog.Info("ping message to client failed", "clientID", c.ID, "connectionId", c.ID)
				return
			}
		}
	}
}
