package client

import (
	"context"
	"net/http"
	"testing"
)

func TestCreateNotificationAddress(t *testing.T) {
	var gotMethod, gotPath string

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 1}`))
	})

	id, err := c.CreateNotificationAddress(context.Background(), NotificationAddress{Address: "a@example.com", NotificationChannelID: "email"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/notification_addresses" {
		t.Errorf("expected POST /notification_addresses, got %s %s", gotMethod, gotPath)
	}
	if id != 1 {
		t.Errorf("expected id 1, got %d", id)
	}
}

func TestReadNotificationAddress(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/notification_addresses/1" {
			t.Errorf("expected path /notification_addresses/1, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"address": "a@example.com", "notification_channel_id": "email"}`))
	})

	address, err := c.ReadNotificationAddress(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if address.Address != "a@example.com" || address.NotificationChannelID != "email" {
		t.Errorf("response didn't unmarshal correctly: %+v", address)
	}
}

func TestReadNotificationAddressNotFound(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{}`))
	})

	_, err := c.ReadNotificationAddress(context.Background(), 999)
	if _, ok := err.(ResourceNotFound); !ok {
		t.Fatalf("expected ResourceNotFound, got %T: %v", err, err)
	}
}

func TestUpdateNotificationAddress(t *testing.T) {
	var gotMethod, gotPath string

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.UpdateNotificationAddress(context.Background(), 1, NotificationAddress{Address: "b@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPut || gotPath != "/notification_addresses/1" {
		t.Errorf("expected PUT /notification_addresses/1, got %s %s", gotMethod, gotPath)
	}
}

func TestDeleteNotificationAddress(t *testing.T) {
	var gotMethod, gotPath string

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.DeleteNotificationAddress(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/notification_addresses/1" {
		t.Errorf("expected DELETE /notification_addresses/1, got %s %s", gotMethod, gotPath)
	}
}
