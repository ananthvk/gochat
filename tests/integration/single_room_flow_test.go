package integration

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ananthvk/gochat/internal/testutils"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func TestWebsocketMessageFlowSingleRoom(t *testing.T) {
	app, srv, cancel := testutils.NewTestServerWithCancel(t)
	defer srv.Close()
	defer cancel()

	room := app.RealtimeService.CreateRoom("test-room-a")

	client1, client1Id := createWSConnection(t, srv)
	client2, client2Id := createWSConnection(t, srv)

	defer client1.Close()
	defer client2.Close()

	// Make both the clients join the same room
	err := app.RealtimeService.JoinRoom(uuid.MustParse(client1Id), room.Id)
	if err != nil {
		t.Errorf("unable to join room clientId=%q roomId=%q", client1Id, room.Id)
	}

	err = app.RealtimeService.JoinRoom(uuid.MustParse(client2Id), room.Id)
	if err != nil {
		t.Errorf("unable to join room clientId=%q roomId=%q", client1Id, room.Id)
	}
	message := "Hello, this is a test message from client 1"
	// Send a message
	err = client1.WriteJSON(map[string]any{
		"type": "chat_message",
		"payload": map[string]any{
			"room_id": room.Id.String(),
			"message": message,
		},
	})
	if err != nil {
		t.Errorf("error while sending ws message %v", err)
	}

	// Receive the message
	resp := make(map[string]any)
	err = client2.ReadJSON(&resp)

	if err != nil {
		t.Errorf("error while reading ws message %v", err)
	}

	msg := resp["payload"].(map[string]any)["message"]
	if msg != message {
		t.Errorf("want %q, got %q, response: %v", message, msg, resp)
	}
}

func createWSConnection(t testing.TB, srv *httptest.Server) (*websocket.Conn, string) {
	t.Helper()
	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/api/v1/realtime/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to websocket: %v", err)
	}

	// Read the initial response containing clientId
	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read websocket message: %v", err)
	}

	var clientResp struct {
		Type    string `json:"type"`
		Payload struct {
			ID string `json:"id"`
		} `json:"payload"`
	}
	if err := json.Unmarshal(message, &clientResp); err != nil {
		t.Fatalf("Failed to unmarshal client response: %v", err)
	}

	if clientResp.Type != "welcome" {
		t.Fatalf("Expected welcome message, got: %s", clientResp.Type)
	}

	clientId := clientResp.Payload.ID
	if clientId == "" {
		t.Fatalf("Response does not contain valid clientId field")
	}
	return conn, clientId
}
