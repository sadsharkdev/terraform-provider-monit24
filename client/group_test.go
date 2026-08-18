package client

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCreateGroup(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody Group

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 7}`))
	})

	assigned := map[string][]int{"default": {1, 2}}

	id, err := c.CreateGroup(context.Background(), Group{
		Name:              "test group",
		AssignedSensorIDs: &assigned,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Errorf("expected method POST, got %s", gotMethod)
	}
	if gotPath != "/groups" {
		t.Errorf("expected path /groups, got %s", gotPath)
	}
	if id != 7 {
		t.Errorf("expected id 7, got %d", id)
	}
	if gotBody.Name != "test group" {
		t.Errorf("expected name=\"test group\", got %q", gotBody.Name)
	}
	if gotBody.AssignedSensorIDs == nil || len((*gotBody.AssignedSensorIDs)["default"]) != 2 {
		t.Errorf("expected assigned_sensor_ids.default=[1 2] in request body, got %v", gotBody.AssignedSensorIDs)
	}
}

// A non-nil, empty AssignedSensorIDs must round-trip as "assigned_sensor_ids":{}
// on the wire, not be omitted — this is what lets monit24_group's
// groupFromResourceData actually clear a previously-set value server-side.
func TestCreateGroupSendsEmptyAssignedSensorIDsAsEmptyObject(t *testing.T) {
	var gotRawBody map[string]interface{}

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotRawBody)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 1}`))
	})

	empty := map[string][]int{}

	_, err := c.CreateGroup(context.Background(), Group{
		Name:              "test group",
		AssignedSensorIDs: &empty,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	raw, present := gotRawBody["assigned_sensor_ids"]
	if !present {
		t.Fatal("expected assigned_sensor_ids to be present in the request body, but it was omitted")
	}
	if m, ok := raw.(map[string]interface{}); !ok || len(m) != 0 {
		t.Errorf("expected assigned_sensor_ids to be an empty object, got %v (%T)", raw, raw)
	}
}

func TestCreateGroupOmitsAssignedSensorIDsWhenNil(t *testing.T) {
	var gotRawBody map[string]interface{}

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotRawBody)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 1}`))
	})

	_, err := c.CreateGroup(context.Background(), Group{Name: "test group"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, present := gotRawBody["assigned_sensor_ids"]; present {
		t.Errorf("expected assigned_sensor_ids to be omitted when nil, got %v", gotRawBody["assigned_sensor_ids"])
	}
}

func TestReadGroup(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/groups/1" {
			t.Errorf("expected path /groups/1, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"name": "test group",
			"is_default": true,
			"periodic_daily_reports": false,
			"assigned_sensor_ids": {"default": [1, 2]}
		}`))
	})

	group, err := c.ReadGroup(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if group.Name != "test group" {
		t.Errorf("expected name=\"test group\", got %q", group.Name)
	}
	if group.IsDefault == nil || !*group.IsDefault {
		t.Errorf("expected is_default=true, got %v", group.IsDefault)
	}
	if group.PeriodicDailyReports == nil || *group.PeriodicDailyReports {
		t.Errorf("expected periodic_daily_reports=false, got %v", group.PeriodicDailyReports)
	}
	if group.AssignedSensorIDs == nil || len((*group.AssignedSensorIDs)["default"]) != 2 {
		t.Errorf("expected assigned_sensor_ids.default=[1 2], got %v", group.AssignedSensorIDs)
	}
}

func TestReadGroupNotFound(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{}`))
	})

	_, err := c.ReadGroup(context.Background(), 999)
	if _, ok := err.(ResourceNotFound); !ok {
		t.Fatalf("expected ResourceNotFound, got %T: %v", err, err)
	}
}

func TestUpdateGroup(t *testing.T) {
	var gotMethod, gotPath string

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.UpdateGroup(context.Background(), 1, Group{Name: "renamed"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotMethod != http.MethodPut {
		t.Errorf("expected method PUT, got %s", gotMethod)
	}
	if gotPath != "/groups/1" {
		t.Errorf("expected path /groups/1, got %s", gotPath)
	}
}

func TestDeleteGroup(t *testing.T) {
	var gotMethod, gotPath string

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.DeleteGroup(context.Background(), 1)
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
