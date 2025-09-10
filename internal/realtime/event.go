package realtime

import "encoding/json"

type event any

// Hub level events

type hubConnectionRegistered struct {
	Client *client
}

type hubConnectionUnregistered struct {
	Client *client
}

type hubDataReceived struct {
	Client *client
	Packet wsDataPacket
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
	Message string `json:"message"`
}
