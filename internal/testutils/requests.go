package testutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ananthvk/gochat/internal/helpers"
)

func CheckStatusCode(t *testing.T, r *http.Response, s int) {
	t.Helper()
	if r.StatusCode != s {
		t.Errorf("wanted status code %q, got %q", s, r.StatusCode)
	}
}

// UnmarshalJSONResponse reads the response body from the http response, and stores the JSON value
// in the value pointed by v.
func UnmarshalJSONResponse(t *testing.T, r *http.Response, v any) {
	defer r.Body.Close()
	t.Helper()
	err := helpers.ParseJSON(r.Body, v, true)

	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
}

// MakePostRequest creates a POST request with JSON body and returns the response
func MakePostRequest(t *testing.T, server *httptest.Server, path string, body any) *http.Response {
	t.Helper()

	jsonData, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	reqBody := bytes.NewReader(jsonData)

	resp, err := http.Post(server.URL+path, "application/json", reqBody)
	if err != nil {
		t.Fatalf("Failed to make POST request: %v", err)
	}

	return resp
}
