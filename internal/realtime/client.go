package realtime

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Maximum number of events that can remain unsent
const maxClientOutgoingSize = 100

// Client represents a websocket connection to the server.
// It stores additional fields such as ConnectedAt,
// and an UUID to uniquely identify this connection
// Outgoing is a buffered channel, and is used by the Hub to send Events to client
type Client struct {
	ID          uuid.UUID
	Connection  *websocket.Conn
	ConnectedAt time.Time
	Outgoing    chan []byte
	Hub         *Hub
}

// NewClient creates a client from the passed websocket connection.
// It sets an unique ID to this connection, and creates the outgoing channel
func NewClient(connection *websocket.Conn, hub *Hub) *Client {
	return &Client{
		ID:          uuid.New(),
		Connection:  connection,
		ConnectedAt: time.Now().UTC(),
		Outgoing:    make(chan []byte, maxClientOutgoingSize),
		Hub:         hub,
	}
}

// ReaderLoop must be run in a separate goroutine. This function runs until the connection is terminated.
// It reads events from the client, and passes those events to the Hub for further processing
func (c *Client) ReaderLoop() {
	defer func() {
		c.Hub.Control <- UnregisterConnectionEvent{Client: c}
		err := c.Connection.Close()
		if err != nil {
			slog.Error("error while closing websocket", "error", err)
		}
		slog.Info("closed websocket", "clientId", c.ID)
	}()

	for {
		messageType, p, err := c.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("websocket read failed", "clientId", c.ID, "error", err)
			}
			return
		}

		// Send the event to the hub for further processing
		// If the Hub Events is full, drop the message, so that the client retransmits it again
		select {
		case c.Hub.Events <- p:
			slog.Info("message enqueued to hub", "clientId", c.ID, "messageType", messageType, "size", len(p))
		default:
			slog.Warn("hub events channel full, dropped message", "clientId", c.ID, "messageType", messageType, "size", len(p))
		}
	}
}
