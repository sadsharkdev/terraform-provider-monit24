package client

import (
	"context"
	"net/http"
	"testing"
)

func TestCreatePeriodicReportAddress(t *testing.T) {
	var gotMethod, gotPath string

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 1}`))
	})

	id, err := c.CreatePeriodicReportAddress(context.Background(), PeriodicReportAddress{Address: "a@example.com", ReportFrequency: "daily"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPost || gotPath != "/periodic_report_addresses" {
		t.Errorf("expected POST /periodic_report_addresses, got %s %s", gotMethod, gotPath)
	}
	if id != 1 {
		t.Errorf("expected id 1, got %d", id)
	}
}

func TestReadPeriodicReportAddress(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/periodic_report_addresses/1" {
			t.Errorf("expected path /periodic_report_addresses/1, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"address": "a@example.com", "report_frequency": "weekly"}`))
	})

	address, err := c.ReadPeriodicReportAddress(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if address.Address != "a@example.com" || address.ReportFrequency != "weekly" {
		t.Errorf("response didn't unmarshal correctly: %+v", address)
	}
}

func TestReadPeriodicReportAddressNotFound(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{}`))
	})

	_, err := c.ReadPeriodicReportAddress(context.Background(), 999)
	if _, ok := err.(ResourceNotFound); !ok {
		t.Fatalf("expected ResourceNotFound, got %T: %v", err, err)
	}
}

func TestUpdatePeriodicReportAddress(t *testing.T) {
	var gotMethod, gotPath string

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.UpdatePeriodicReportAddress(context.Background(), 1, PeriodicReportAddress{Address: "b@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPut || gotPath != "/periodic_report_addresses/1" {
		t.Errorf("expected PUT /periodic_report_addresses/1, got %s %s", gotMethod, gotPath)
	}
}

func TestDeletePeriodicReportAddress(t *testing.T) {
	var gotMethod, gotPath string

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.DeletePeriodicReportAddress(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/periodic_report_addresses/1" {
		t.Errorf("expected DELETE /periodic_report_addresses/1, got %s %s", gotMethod, gotPath)
	}
}
