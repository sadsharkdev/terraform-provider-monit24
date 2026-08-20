package monit24

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/monit24/terraform-provider-monit24/client"
)

func TestAccPeriodicReportAddress(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPeriodicReportAddressConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_periodic_report_address.test", "address", "reports@example.com"),
					resource.TestCheckResourceAttr("monit24_periodic_report_address.test", "report_frequency", "daily"),
				),
			},
			{
				Config: testAccPeriodicReportAddressConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_periodic_report_address.test", "address", "reports@example.com"),
					resource.TestCheckResourceAttr("monit24_periodic_report_address.test", "report_frequency", "weekly"),
				),
			},
		},
	})
}

const testAccPeriodicReportAddressConfig = `
resource "monit24_periodic_report_address" "test" {
  address           = "reports@example.com"
  report_frequency  = "daily"
  group_id          = monit24_group.test.id
}

resource "monit24_group" "test" {
  name = "test group"
}
`

const testAccPeriodicReportAddressConfigUpdated = `
resource "monit24_periodic_report_address" "test" {
  address           = "reports@example.com"
  report_frequency  = "weekly"
  group_id          = monit24_group.test.id
}

resource "monit24_group" "test" {
  name = "test group"
}
`

func TestPeriodicReportAddressFromResourceData(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourcePeriodicReportAddress().Schema, map[string]interface{}{
		"address":          "reports@example.com",
		"report_frequency": "weekly",
		"group_id":         5,
	})

	address := periodicReportAddressFromResourceData(d, client.Client{})

	if address.Address != "reports@example.com" || address.ReportFrequency != "weekly" {
		t.Errorf("expected address/report_frequency to round-trip, got %+v", address)
	}
	if address.GroupID != 5 {
		t.Errorf("expected group_id=5, got %d", address.GroupID)
	}
}

func TestPeriodicReportAddressCreateUpdateDeleteLifecycle(t *testing.T) {
	var stored client.PeriodicReportAddress

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/periodic_report_addresses":
			json.NewDecoder(r.Body).Decode(&stored)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]int{"id": 1})
		case r.Method == http.MethodGet && r.URL.Path == "/periodic_report_addresses/1":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(stored)
		case r.Method == http.MethodPut && r.URL.Path == "/periodic_report_addresses/1":
			json.NewDecoder(r.Body).Decode(&stored)
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodDelete && r.URL.Path == "/periodic_report_addresses/1":
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	c := client.NewTestClient(server.URL)

	d := schema.TestResourceDataRaw(t, resourcePeriodicReportAddress().Schema, map[string]interface{}{
		"address":          "reports@example.com",
		"report_frequency": "daily",
	})

	if diags := resourcePeriodicReportAddressCreate(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error creating periodic_report_address: %v", diags)
	}
	if d.Id() != "1" {
		t.Fatalf("expected id \"1\", got %q", d.Id())
	}
	if stored.ReportFrequency != "daily" {
		t.Errorf("expected create request to carry report_frequency, server stored %+v", stored)
	}

	if diags := resourcePeriodicReportAddressRead(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error reading periodic_report_address: %v", diags)
	}
	if d.Get("address").(string) != "reports@example.com" {
		t.Errorf("expected address to round-trip through read, got %q", d.Get("address"))
	}

	if err := d.Set("report_frequency", "weekly"); err != nil {
		t.Fatalf("unexpected error setting report_frequency: %v", err)
	}
	if diags := resourcePeriodicReportAddressUpdate(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error updating periodic_report_address: %v", diags)
	}
	if stored.ReportFrequency != "weekly" {
		t.Errorf("expected update request to carry the new report_frequency, server stored %+v", stored)
	}

	if diags := resourcePeriodicReportAddressDelete(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error deleting periodic_report_address: %v", diags)
	}
}
