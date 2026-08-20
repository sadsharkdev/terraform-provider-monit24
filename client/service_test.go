package client

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCreateService(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody Service

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 42}`))
	})

	req := Service{
		TypeID:    "https",
		Name:      "example",
		Address:   "example.com",
		Interval:  600,
		SensorIDs: &[]int{1, 2},
		StepNames: &[]string{"step 1", "step 2"},
		ExtendedSettings: &map[string]interface{}{
			"http_method": "POST",
			"port":        443,
		},
	}

	id, err := c.CreateService(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Errorf("expected method POST, got %s", gotMethod)
	}
	if gotPath != "/services" {
		t.Errorf("expected path /services, got %s", gotPath)
	}
	if id != 42 {
		t.Errorf("expected id 42, got %d", id)
	}

	if gotBody.TypeID != "https" || gotBody.Address != "example.com" {
		t.Errorf("request body didn't round-trip correctly: %+v", gotBody)
	}
	if gotBody.SensorIDs == nil || len(*gotBody.SensorIDs) != 2 {
		t.Errorf("expected sensor_ids [1 2] in request body, got %v", gotBody.SensorIDs)
	}
	if gotBody.StepNames == nil || len(*gotBody.StepNames) != 2 {
		t.Errorf("expected step_names in request body, got %v", gotBody.StepNames)
	}
	if gotBody.ExtendedSettings == nil || (*gotBody.ExtendedSettings)["http_method"] != "POST" {
		t.Errorf("expected extended_settings.http_method=POST in request body, got %v", gotBody.ExtendedSettings)
	}
}

func TestReadService(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected method GET, got %s", r.Method)
		}
		if r.URL.Path != "/services/123" {
			t.Errorf("expected path /services/123, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"type_id": "https",
			"name": "example",
			"address": "example.com",
			"interval": 600,
			"is_active": true,
			"is_archived": false,
			"sensor_ids": [1, 2],
			"extended_settings": {"http_method": "GET"}
		}`))
	})

	service, err := c.ReadService(context.Background(), 123)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if service.TypeID != "https" || service.Address != "example.com" {
		t.Errorf("response didn't unmarshal correctly: %+v", service)
	}
	if service.IsActive == nil || !*service.IsActive {
		t.Errorf("expected is_active=true, got %v", service.IsActive)
	}
	if service.IsArchived == nil || *service.IsArchived {
		t.Errorf("expected is_archived=false, got %v", service.IsArchived)
	}
	if service.SensorIDs == nil || len(*service.SensorIDs) != 2 {
		t.Errorf("expected sensor_ids [1 2], got %v", service.SensorIDs)
	}
}

func TestReadServiceNotFound(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{}`))
	})

	_, err := c.ReadService(context.Background(), 999)
	if _, ok := err.(ResourceNotFound); !ok {
		t.Fatalf("expected ResourceNotFound, got %T: %v", err, err)
	}
}

func TestUpdateService(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody Service

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.UpdateService(context.Background(), 123, Service{Name: "renamed"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotMethod != http.MethodPut {
		t.Errorf("expected method PUT, got %s", gotMethod)
	}
	if gotPath != "/services/123" {
		t.Errorf("expected path /services/123, got %s", gotPath)
	}
	if gotBody.Name != "renamed" {
		t.Errorf("expected request body name=renamed, got %q", gotBody.Name)
	}
}

func TestDeleteService(t *testing.T) {
	var gotMethod, gotPath string

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.DeleteService(context.Background(), 123)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotMethod != http.MethodDelete {
		t.Errorf("expected method DELETE, got %s", gotMethod)
	}
	if gotPath != "/services/123" {
		t.Errorf("expected path /services/123, got %s", gotPath)
	}
}
