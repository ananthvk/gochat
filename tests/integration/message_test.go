package integration

import (
	"context"
	"net/http"
	"testing"

	"github.com/ananthvk/gochat/internal/auth"
	"github.com/ananthvk/gochat/internal/testutils"
	"github.com/oklog/ulid/v2"
)

func TestMessage(t *testing.T) {
	app, srv, _, cancel := testutils.NewTestServerWithDatabaseAndCancel(t)
	defer srv.Close()
	defer cancel()

	// Create a group first for message testing
	createGroupResp := testutils.MakePostRequest(t, srv, "/api/v1/group", map[string]any{
		"name":        "Message Test Group",
		"description": "Group for testing message functionality",
	})
	testutils.CheckStatusCode(t, createGroupResp, http.StatusCreated)

	createGroupData := map[string]any{}
	testutils.UnmarshalJSONResponse(t, createGroupResp, &createGroupData)
	groupId := createGroupData["id"].(string)

	t.Run("TestMessageCreation", func(t *testing.T) {
		resp := testutils.MakePostRequest(t, srv, "/api/v1/group/"+groupId+"/message", map[string]any{
			"content": "Hello, this is a test message!",
			"type":    "text",
		})
		testutils.CheckStatusCode(t, resp, http.StatusCreated)

		respData := map[string]any{}
		testutils.UnmarshalJSONResponse(t, resp, &respData)

		messageId, ok := respData["id"].(string)
		if !ok {
			t.Fatalf("Response does not contain valid id field")
		}
		dbMsg, _ := app.MessageService.GetOne(context.Background(), ulid.MustParse(messageId), ulid.MustParse(groupId), ulid.MustParse(auth.HardcodedUserId))
		if messageId != ulid.ULID(dbMsg.ID).String() {
			t.Errorf("expected %q, got %q", messageId, dbMsg.ID)
		}
	})

	t.Run("TestMessageFetch", func(t *testing.T) {
		// Create a message first
		createResp := testutils.MakePostRequest(t, srv, "/api/v1/group/"+groupId+"/message", map[string]any{
			"content": "Fetch test message",
			"type":    "text",
		})
		testutils.CheckStatusCode(t, createResp, http.StatusCreated)

		createData := map[string]any{}
		testutils.UnmarshalJSONResponse(t, createResp, &createData)
		messageId := createData["id"].(string)

		// Fetch the message
		resp := testutils.MakeGetRequest(t, srv, "/api/v1/group/"+groupId+"/message/"+messageId)
		testutils.CheckStatusCode(t, resp, http.StatusOK)

		respData := map[string]any{}
		testutils.UnmarshalJSONResponse(t, resp, &respData)

		if respData["content"] != "Fetch test message" {
			t.Errorf("expected content %q, got %q", "Fetch test message", respData["content"])
		}
	})

	t.Run("TestMessageDelete", func(t *testing.T) {
		// Create a message first
		createResp := testutils.MakePostRequest(t, srv, "/api/v1/group/"+groupId+"/message", map[string]any{
			"content": "Message to delete",
			"type":    "text",
		})
		testutils.CheckStatusCode(t, createResp, http.StatusCreated)

		createData := map[string]any{}
		testutils.UnmarshalJSONResponse(t, createResp, &createData)
		messageId := createData["id"].(string)

		// Delete the message
		resp := testutils.MakeDeleteRequest(t, srv, "/api/v1/group/"+groupId+"/message/"+messageId)
		testutils.CheckStatusCode(t, resp, http.StatusOK)

		// Verify deletion by attempting to fetch
		fetchResp := testutils.MakeGetRequest(t, srv, "/api/v1/group/"+groupId+"/message/"+messageId)
		testutils.CheckStatusCode(t, fetchResp, http.StatusNotFound)
	})

	t.Run("TestMessageListByGroup", func(t *testing.T) {
		// Create multiple messages in the same group
		for i := 0; i < 3; i++ {
			createResp := testutils.MakePostRequest(t, srv, "/api/v1/group/"+groupId+"/message", map[string]any{
				"content": "List test message " + string(rune(i+'1')),
				"type":    "text",
			})
			testutils.CheckStatusCode(t, createResp, http.StatusCreated)
		}

		// List messages by group
		resp := testutils.MakeGetRequest(t, srv, "/api/v1/group/"+groupId+"/message")
		testutils.CheckStatusCode(t, resp, http.StatusOK)

		respData := map[string]any{}
		testutils.UnmarshalJSONResponse(t, resp, &respData)

		messages, ok := respData["messages"].([]any)
		if !ok {
			t.Fatalf("Response does not contain valid messages field")
		}

		if len(messages) < 3 {
			t.Errorf("expected at least 3 messages, got %d", len(messages))
		}
	})
}
