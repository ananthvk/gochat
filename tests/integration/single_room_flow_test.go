package integration

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ananthvk/gochat/internal/testutils"
	"github.com/gorilla/websocket"
)

func TestWebsocketMessageFlowSingleRoom(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	_, router := testutils.NewTestServer(t, ctx)
	srv := httptest.NewServer(router)
	defer srv.Close()
	defer cancel()

	reqBody := strings.NewReader(`{"name": "test-room"}`)
	resp, err := http.Post(srv.URL+"/api/v1/realtime/room", "application/json", reqBody)
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201 response, got %v", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var respData map[string]any
	if err := json.Unmarshal(body, &respData); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	_, ok := respData["id"].(string)
	if !ok {
		t.Fatalf("Response does not contain valid id field")
	}

	defer resp.Body.Close()

	conn1 := createWSConnection(t, srv.URL)
	defer conn1.Close()
	conn2 := createWSConnection(t, srv.URL)
	defer conn2.Close()

}

func createWSConnection(t testing.TB, URL string) *websocket.Conn {
	t.Helper()
	wsURL := "ws" + strings.TrimPrefix(URL, "http") + "/api/v1/realtime/ws"
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
	return conn
}
