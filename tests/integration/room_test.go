package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/ananthvk/gochat/internal/testutils"
)

// TestRoomCreation tests whether the api endpoint to create a room works correctly
func TestRoomCreation(t *testing.T) {
	app, srv, cancel := testutils.NewTestServerWithCancel(t)
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

	roomId, ok := respData["id"].(string)
	if !ok {
		t.Fatalf("Response does not contain valid id field")
	}

	if app.RealtimeService.GetRoomByName("test-room").Id.String() != roomId {
		t.Fatalf("returned id does not match with actually created one")
	}
	defer resp.Body.Close()
}

// TestRoomByName tests whether the api endpoint to fetch a room by name works correctly
func TestRoomByName(t *testing.T) {
	app, srv, cancel := testutils.NewTestServerWithCancel(t)
	defer srv.Close()
	defer cancel()

	roomName := "test-room-1"

	room := app.RealtimeService.CreateRoom(roomName)

	resp, err := http.Get(srv.URL + "/api/v1/realtime/room/by-name/" + url.PathEscape(roomName))
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 response, got %v", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var respData map[string]any
	if err := json.Unmarshal(body, &respData); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	roomId, ok := respData["id"].(string)
	if !ok {
		t.Fatalf("Response does not contain valid id field")
	}

	if roomId != room.Id.String() {
		t.Errorf("expected %q, got %q", room.Id, roomId)
	}
}
