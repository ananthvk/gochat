package integration

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/ananthvk/gochat/internal/testutils"
)

// TestRoomCreation tests whether the api endpoint to create a room works correctly
func TestRoomCreation(t *testing.T) {
	app, srv, cancel := testutils.NewTestServerWithCancel(t)
	defer srv.Close()
	defer cancel()

	resp := testutils.MakePostRequest(
		t, srv, "/api/v1/realtime/room", map[string]any{
			"name": "test-room",
		})
	testutils.CheckStatusCode(t, resp, http.StatusCreated)

	respData := map[string]any{}
	testutils.UnmarshalJSONResponse(t, resp, &respData)

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
		t.Fatalf("Failed to make GET request: %v", err)
	}
	testutils.CheckStatusCode(t, resp, http.StatusOK)

	respData := map[string]any{}
	testutils.UnmarshalJSONResponse(t, resp, &respData)

	roomId, ok := respData["id"].(string)
	if !ok {
		t.Fatalf("Response does not contain valid id field")
	}

	if roomId != room.Id.String() {
		t.Errorf("expected %q, got %q", room.Id, roomId)
	}
}
