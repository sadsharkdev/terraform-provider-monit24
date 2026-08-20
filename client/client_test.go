package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
)

func TestAuthorizationHeaderValue(t *testing.T) {
	tests := []struct {
		name      string
		basicAuth string
		token     string
		want      string
	}{
		{
			name:      "token takes priority when both are set",
			basicAuth: "dXNlcjpwYXNz",
			token:     "abc123",
			want:      "Bearer abc123",
		},
		{
			name:      "basic auth used when no token",
			basicAuth: "dXNlcjpwYXNz",
			token:     "",
			want:      "Basic dXNlcjpwYXNz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := authorizationHeaderValue(tt.basicAuth, tt.token)
			if got != tt.want {
				t.Errorf("authorizationHeaderValue(%q, %q) = %q, want %q", tt.basicAuth, tt.token, got, tt.want)
			}
		})
	}
}

func TestGetSendsCorrectMethodAndPath(t *testing.T) {
	var gotMethod, gotPath string

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})

	_, err := c.get(context.Background(), "/services/123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotMethod != http.MethodGet {
		t.Errorf("expected method GET, got %s", gotMethod)
	}
	if gotPath != "/services/123" {
		t.Errorf("expected path /services/123, got %s", gotPath)
	}
}

func TestPostEncodesPayloadAndExpectsCreated(t *testing.T) {
	type payload struct {
		Name string `json:"name"`
	}

	var gotMethod string
	var gotBody payload

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 42}`))
	})

	resp, err := c.post(context.Background(), "/groups", payload{Name: "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Errorf("expected method POST, got %s", gotMethod)
	}
	if gotBody.Name != "test" {
		t.Errorf("expected request body name=test, got %q", gotBody.Name)
	}

	var parsed struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(resp, &parsed); err != nil {
		t.Fatalf("unexpected error unmarshaling response: %v", err)
	}
	if parsed.ID != 42 {
		t.Errorf("expected response id=42, got %d", parsed.ID)
	}
}

func TestPutEncodesPayloadAndExpectsNoContent(t *testing.T) {
	type payload struct {
		Name string `json:"name"`
	}

	var gotMethod string
	var gotBody payload

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.put(context.Background(), "/groups/1", payload{Name: "updated"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotMethod != http.MethodPut {
		t.Errorf("expected method PUT, got %s", gotMethod)
	}
	if gotBody.Name != "updated" {
		t.Errorf("expected request body name=updated, got %q", gotBody.Name)
	}
}

func TestDeleteSendsCorrectMethodAndExpectsNoContent(t *testing.T) {
	var gotMethod, gotPath string

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.delete(context.Background(), "/groups/1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotMethod != http.MethodDelete {
		t.Errorf("expected method DELETE, got %s", gotMethod)
	}
	if gotPath != "/groups/1" {
		t.Errorf("expected path /groups/1, got %s", gotPath)
	}
}

func TestRequestReturnsResourceNotFoundOn404(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "not found"}`))
	})

	_, err := c.get(context.Background(), "/groups/999")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	var notFound ResourceNotFound
	if !errors.As(err, &notFound) {
		t.Fatalf("expected client.ResourceNotFound, got %T: %v", err, err)
	}
}

func TestRequestReturnsGenericErrorOnUnexpectedStatus(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "boom"}`))
	})

	_, err := c.get(context.Background(), "/groups/1")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	var notFound ResourceNotFound
	if errors.As(err, &notFound) {
		t.Fatal("expected a generic error, got client.ResourceNotFound")
	}

	if !strings.Contains(err.Error(), "500") || !strings.Contains(err.Error(), "boom") {
		t.Errorf("expected error message to mention status code and body, got: %v", err)
	}
}

func TestRequestSendsAuthorizationHeader(t *testing.T) {
	var gotAuth string

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})
	c.token = "my-token"

	_, err := c.get(context.Background(), "/accounts/my_account")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotAuth != "Bearer my-token" {
		t.Errorf("expected Authorization: Bearer my-token, got %q", gotAuth)
	}
}
