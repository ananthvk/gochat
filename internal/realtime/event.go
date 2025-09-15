package realtime

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type event any

// Hub level events

type hubConnectionRegistered struct {
	ClientId   uuid.UUID
	Connection *websocket.Conn
}

type hubConnectionUnregistered struct {
	ClientId uuid.UUID
}

type hubDataReceived struct {
	ClientId uuid.UUID
	Packet   wsDataPacket
}

type hubRoomCreated struct {
	roomId uuid.UUID
}

type hubRoomJoined struct {
	ClientId uuid.UUID
	roomId   uuid.UUID
}

// Application level events (Data of a hubDataReceived event) contains one of these

// wsDataPacket contains the data received from a websocket. It's a JSON message containing two fields.
// A Type field, and a Payload field. The Type field determines the type of message, and the paylod contains any JSON object
// that can be parsed later depending upon the type.
type wsDataPacket struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// wsChatMessage represents a chat message either sent from the client, or from the server.
// Type is set to "chat_message"
type wsChatMessage struct {
	RoomID  uuid.UUID `json:"room_id"`
	Message string    `json:"message"`
}
