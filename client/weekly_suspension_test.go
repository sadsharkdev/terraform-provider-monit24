package client

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCreateWeeklySuspension(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody WeeklySuspension

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 1}`))
	})

	id, err := c.CreateWeeklySuspension(context.Background(), WeeklySuspension{
		ServiceID:   1,
		StartMinute: MinuteOfWeek{DayOfWeek: 6, Hour: 22, Minute: 0},
		EndMinute:   MinuteOfWeek{DayOfWeek: 7, Hour: 2, Minute: 0},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/weekly_suspensions" {
		t.Errorf("expected POST /weekly_suspensions, got %s %s", gotMethod, gotPath)
	}
	if id != 1 {
		t.Errorf("expected id 1, got %d", id)
	}
	if gotBody.StartMinute.DayOfWeek != 6 || gotBody.StartMinute.Hour != 22 {
		t.Errorf("expected start_minute nested object to round-trip, got %+v", gotBody.StartMinute)
	}
	if gotBody.EndMinute.DayOfWeek != 7 || gotBody.EndMinute.Hour != 2 {
		t.Errorf("expected end_minute nested object to round-trip, got %+v", gotBody.EndMinute)
	}
}

func TestReadWeeklySuspension(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/weekly_suspensions/1" {
			t.Errorf("expected path /weekly_suspensions/1, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"service_id": 1,
			"start_minute": {"day_of_week": 6, "hour": 22, "minute": 0},
			"end_minute": {"day_of_week": 7, "hour": 2, "minute": 0}
		}`))
	})

	suspension, err := c.ReadWeeklySuspension(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if suspension.StartMinute.DayOfWeek != 6 || suspension.EndMinute.DayOfWeek != 7 {
		t.Errorf("response didn't unmarshal correctly: %+v", suspension)
	}
}

func TestReadWeeklySuspensionNotFound(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{}`))
	})

	_, err := c.ReadWeeklySuspension(context.Background(), 999)
	if _, ok := err.(ResourceNotFound); !ok {
		t.Fatalf("expected ResourceNotFound, got %T: %v", err, err)
	}
}

func TestUpdateWeeklySuspension(t *testing.T) {
	var gotMethod, gotPath string

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.UpdateWeeklySuspension(context.Background(), 1, WeeklySuspension{ServiceID: 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPut || gotPath != "/weekly_suspensions/1" {
		t.Errorf("expected PUT /weekly_suspensions/1, got %s %s", gotMethod, gotPath)
	}
}

func TestDeleteWeeklySuspension(t *testing.T) {
	var gotMethod, gotPath string

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.DeleteWeeklySuspension(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/weekly_suspensions/1" {
		t.Errorf("expected DELETE /weekly_suspensions/1, got %s %s", gotMethod, gotPath)
	}
}
