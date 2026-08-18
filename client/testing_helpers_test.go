package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// newTestServer starts an httptest.Server driven by handler and returns a
// Client pointed at it, so client methods can be exercised against a fake
// Monit24 API without a network call or live credentials.
func newTestServer(t *testing.T, handler http.HandlerFunc) Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return NewTestClient(server.URL)
}
