package client

import (
	"context"
	"net/http"
	"testing"
)

func TestCreateSuspension(t *testing.T) {
	var gotMethod, gotPath string

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 1}`))
	})

	id, err := c.CreateSuspension(context.Background(), Suspension{ServiceID: 1, EndTime: "2099-01-01T00:00:00Z"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/suspensions" {
		t.Errorf("expected POST /suspensions, got %s %s", gotMethod, gotPath)
	}
	if id != 1 {
		t.Errorf("expected id 1, got %d", id)
	}
}

func TestReadSuspension(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/suspensions/1" {
			t.Errorf("expected path /suspensions/1, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"service_id": 1, "end_time": "2099-01-01T00:00:00Z"}`))
	})

	suspension, err := c.ReadSuspension(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if suspension.ServiceID != 1 || suspension.EndTime != "2099-01-01T00:00:00Z" {
		t.Errorf("response didn't unmarshal correctly: %+v", suspension)
	}
}

func TestReadSuspensionNotFound(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{}`))
	})

	_, err := c.ReadSuspension(context.Background(), 999)
	if _, ok := err.(ResourceNotFound); !ok {
		t.Fatalf("expected ResourceNotFound, got %T: %v", err, err)
	}
}

func TestUpdateSuspension(t *testing.T) {
	var gotMethod, gotPath string

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.UpdateSuspension(context.Background(), 1, Suspension{ServiceID: 1, EndTime: "2099-06-01T00:00:00Z"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPut || gotPath != "/suspensions/1" {
		t.Errorf("expected PUT /suspensions/1, got %s %s", gotMethod, gotPath)
	}
}

func TestDeleteSuspension(t *testing.T) {
	var gotMethod, gotPath string

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.DeleteSuspension(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/suspensions/1" {
		t.Errorf("expected DELETE /suspensions/1, got %s %s", gotMethod, gotPath)
	}
}
