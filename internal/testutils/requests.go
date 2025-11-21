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
		t.Errorf("wanted status code %d, got %d", s, r.StatusCode)
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

// MakeGetRequest creates a GET request and returns the response
func MakeGetRequest(t *testing.T, server *httptest.Server, path string) *http.Response {
	t.Helper()

	resp, err := http.Get(server.URL + path)
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}

	return resp
}

// MakeDeleteRequest creates a DELETE request and returns the response
func MakeDeleteRequest(t *testing.T, server *httptest.Server, path string) *http.Response {
	t.Helper()

	req, err := http.NewRequest(http.MethodDelete, server.URL+path, nil)
	if err != nil {
		t.Fatalf("Failed to create DELETE request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make DELETE request: %v", err)
	}

	return resp
}

// MakePatchRequest creates a PATCH request with JSON body and returns the response
func MakePatchRequest(t *testing.T, server *httptest.Server, path string, body any) *http.Response {
	t.Helper()

	jsonData, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	reqBody := bytes.NewReader(jsonData)

	req, err := http.NewRequest(http.MethodPatch, server.URL+path, reqBody)
	if err != nil {
		t.Fatalf("Failed to create PATCH request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make PATCH request: %v", err)
	}

	return resp
}

type AuthenticatedRequest struct {
	UserId string
	Token  string
}

// GetAuth should be called on a newly created AuthenticatedRequest instance, it registers a new user using /signup and sets the returned token so that
// authenticated requests can be made
func (a *AuthenticatedRequest) GetAuth(t *testing.T, server *httptest.Server) {
	t.Helper()
	resp := MakePostRequest(t, server, "/api/v1/auth/signup", map[string]any{
		"name":     "A Test User",
		"email":    "test@example.com",
		"username": "test_user",
		"password": "testuser123",
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create user failed, expected 201 response, got %d", resp.StatusCode)
	}
	type User struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		Activated bool   `json:"activated"`
		CreatedAt string `json:"created_at"`
	}
	type Token struct {
		Expiry string `json:"expiry"`
		Token  string `json:"token"`
	}
	type signupResponse struct {
		User  User  `json:"user"`
		Token Token `json:"token"`
	}
	signup := signupResponse{}
	UnmarshalJSONResponse(t, resp, &signup)

	// Set the user id
	a.UserId = signup.User.ID
	a.Token = signup.Token.Token
}

// MakeAuthenticatedPostRequest creates a POST request with JSON body and Authorization header
func (a *AuthenticatedRequest) MakeAuthenticatedPostRequest(t *testing.T, server *httptest.Server, path string, body any) *http.Response {
	t.Helper()

	jsonData, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	reqBody := bytes.NewReader(jsonData)

	req, err := http.NewRequest(http.MethodPost, server.URL+path, reqBody)
	if err != nil {
		t.Fatalf("Failed to create POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make authenticated POST request: %v", err)
	}

	return resp
}

// MakeAuthenticatedGetRequest creates a GET request with Authorization header
func (a *AuthenticatedRequest) MakeAuthenticatedGetRequest(t *testing.T, server *httptest.Server, path string) *http.Response {
	t.Helper()

	req, err := http.NewRequest(http.MethodGet, server.URL+path, nil)
	if err != nil {
		t.Fatalf("Failed to create GET request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+a.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make authenticated GET request: %v", err)
	}

	return resp
}

// MakeAuthenticatedDeleteRequest creates a DELETE request with Authorization header
func (a *AuthenticatedRequest) MakeAuthenticatedDeleteRequest(t *testing.T, server *httptest.Server, path string) *http.Response {
	t.Helper()

	req, err := http.NewRequest(http.MethodDelete, server.URL+path, nil)
	if err != nil {
		t.Fatalf("Failed to create DELETE request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+a.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make authenticated DELETE request: %v", err)
	}

	return resp
}

// MakeAuthenticatedPatchRequest creates a PATCH request with JSON body and Authorization header
func (a *AuthenticatedRequest) MakeAuthenticatedPatchRequest(t *testing.T, server *httptest.Server, path string, body any) *http.Response {
	t.Helper()

	jsonData, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	reqBody := bytes.NewReader(jsonData)

	req, err := http.NewRequest(http.MethodPatch, server.URL+path, reqBody)
	if err != nil {
		t.Fatalf("Failed to create PATCH request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make authenticated PATCH request: %v", err)
	}

	return resp
}
