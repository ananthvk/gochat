package integration

import (
	"context"
	"net/http"
	"testing"

	"github.com/ananthvk/gochat/internal/testutils"
	"github.com/oklog/ulid/v2"
)

func TestGroup(t *testing.T) {
	app, srv, _, cancel := testutils.NewTestServerWithDatabaseAndCancel(t)
	defer srv.Close()
	defer cancel()
	req := testutils.AuthenticatedRequest{}
	req.GetAuth(t, srv)

	// TestGroupCreation tests whether the api endpoint to fetch a create a room works correctly
	t.Run("TestGroupCreation", func(t *testing.T) {
		resp := req.MakeAuthenticatedPostRequest(t, srv, "/api/v1/group", map[string]any{
			"name":        "A Test Group",
			"description": "A new group to test out the features",
		})
		testutils.CheckStatusCode(t, resp, http.StatusCreated)

		respData := map[string]any{}
		testutils.UnmarshalJSONResponse(t, resp, &respData)

		groupId, ok := respData["id"].(string)
		if !ok {
			t.Fatalf("Response does not contain valid id field")
		}
		dbGrp, _ := app.GroupService.GetOne(context.Background(), ulid.MustParse(groupId), ulid.MustParse(req.UserId))
		if groupId != ulid.ULID(dbGrp.ID).String() {
			t.Errorf("expected %q, got %q", groupId, dbGrp.ID)
		}
	})

	t.Run("TestGroupFetch", func(t *testing.T) {
		// Create a group first
		createResp := req.MakeAuthenticatedPostRequest(t, srv, "/api/v1/group", map[string]any{
			"name":        "Fetch Test Group",
			"description": "Group for testing fetch functionality",
		})
		testutils.CheckStatusCode(t, createResp, http.StatusCreated)

		createData := map[string]any{}
		testutils.UnmarshalJSONResponse(t, createResp, &createData)
		groupId := createData["id"].(string)

		// Fetch the group
		resp := req.MakeAuthenticatedGetRequest(t, srv, "/api/v1/group/"+groupId)
		testutils.CheckStatusCode(t, resp, http.StatusOK)

		respData := map[string]any{}
		testutils.UnmarshalJSONResponse(t, resp, &respData)

		if respData["name"] != "Fetch Test Group" {
			t.Errorf("expected name %q, got %q", "Fetch Test Group", respData["name"])
		}
	})

	t.Run("TestGroupUpdate", func(t *testing.T) {
		// Create a group first
		createResp := req.MakeAuthenticatedPostRequest(t, srv, "/api/v1/group", map[string]any{
			"name":        "Original Name",
			"description": "Original description",
		})
		testutils.CheckStatusCode(t, createResp, http.StatusCreated)

		createData := map[string]any{}
		testutils.UnmarshalJSONResponse(t, createResp, &createData)
		groupId := createData["id"].(string)

		// Update the group
		resp := req.MakeAuthenticatedPatchRequest(t, srv, "/api/v1/group/"+groupId, map[string]any{
			"name": "Updated Name",
		})
		testutils.CheckStatusCode(t, resp, http.StatusOK)

		// Verify the update
		dbGrp, _ := app.GroupService.GetOne(context.Background(), ulid.MustParse(groupId), ulid.MustParse(req.UserId))
		if dbGrp.Name != "Updated Name" {
			t.Errorf("expected name %q, got %q", "Updated Name", dbGrp.Name)
		}
	})

	t.Run("TestGroupDelete", func(t *testing.T) {
		// Create a group first
		createResp := req.MakeAuthenticatedPostRequest(t, srv, "/api/v1/group", map[string]any{
			"name":        "Group to Delete",
			"description": "This group will be deleted",
		})
		testutils.CheckStatusCode(t, createResp, http.StatusCreated)

		createData := map[string]any{}
		testutils.UnmarshalJSONResponse(t, createResp, &createData)
		groupId := createData["id"].(string)

		// Delete the group
		resp := req.MakeAuthenticatedDeleteRequest(t, srv, "/api/v1/group/"+groupId)
		testutils.CheckStatusCode(t, resp, http.StatusOK)

		// Verify deletion by attempting to fetch
		fetchResp := req.MakeAuthenticatedGetRequest(t, srv, "/api/v1/group/"+groupId)
		testutils.CheckStatusCode(t, fetchResp, http.StatusNotFound)
	})

	t.Run("TestGroupListAll", func(t *testing.T) {
		// Create multiple groups
		for i := 0; i < 3; i++ {
			createResp := req.MakeAuthenticatedPostRequest(t, srv, "/api/v1/group", map[string]any{
				"name":        "List Test Group " + string(rune(i+'1')),
				"description": "Group for testing list functionality",
			})
			testutils.CheckStatusCode(t, createResp, http.StatusCreated)
		}

		// List all groups
		resp := req.MakeAuthenticatedGetRequest(t, srv, "/api/v1/group")
		testutils.CheckStatusCode(t, resp, http.StatusOK)

		respData := map[string]any{}

		testutils.UnmarshalJSONResponse(t, resp, &respData)

		groups, ok := respData["groups"].([]any)
		if !ok {
			t.Fatalf("Response does not contain valid groups field")
		}

		if len(groups) < 3 {
			t.Errorf("expected at least 3 groups, got %d", len(groups))
		}
	})
}
