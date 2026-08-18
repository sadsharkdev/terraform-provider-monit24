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

func TestAccSuspension(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSuspensionConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_suspension.test", "end_time", "2099-01-01T00:00:00Z"),
					resource.TestCheckResourceAttr("monit24_suspension.test", "only_notifications", "false"),
				),
			},
			{
				Config: testAccSuspensionConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_suspension.test", "end_time", "2099-06-01T00:00:00Z"),
					resource.TestCheckResourceAttr("monit24_suspension.test", "only_notifications", "true"),
					resource.TestCheckResourceAttr("monit24_suspension.test", "description", "planned maintenance"),
				),
			},
		},
	})
}

const testAccSuspensionConfig = `
resource "monit24_suspension" "test" {
  service_id = monit24_service.test.id
  end_time   = "2099-01-01T00:00:00Z"
}

resource "monit24_service" "test" {
  type_id     = "https"
  name        = "https example.com"
  address     = "example.com"
  group_id    = monit24_group.test.id
}

resource "monit24_group" "test" {
  name = "test group"
}
`

const testAccSuspensionConfigUpdated = `
resource "monit24_suspension" "test" {
  service_id          = monit24_service.test.id
  end_time             = "2099-06-01T00:00:00Z"
  only_notifications   = true
  description          = "planned maintenance"
}

resource "monit24_service" "test" {
  type_id     = "https"
  name        = "https example.com"
  address     = "example.com"
  group_id    = monit24_group.test.id
}

resource "monit24_group" "test" {
  name = "test group"
}
`

func TestSuspensionFromResourceData(t *testing.T) {
	t.Run("required fields, start_time and description omitted", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceSuspension().Schema, map[string]interface{}{
			"service_id": 1,
			"end_time":   "2099-01-01T00:00:00Z",
		})

		suspension := suspensionFromResourceData(d)

		if suspension.ServiceID != 1 || suspension.EndTime != "2099-01-01T00:00:00Z" {
			t.Errorf("expected service_id/end_time to round-trip, got %+v", suspension)
		}
		if suspension.StartTime != nil {
			t.Errorf("expected start_time to stay nil when not configured, got %v", *suspension.StartTime)
		}
		if suspension.Description != nil {
			t.Errorf("expected description to stay nil when not configured, got %v", *suspension.Description)
		}
		if suspension.OnlyNotifications == nil || *suspension.OnlyNotifications {
			t.Errorf("expected only_notifications default false, got %v", suspension.OnlyNotifications)
		}
	})

	t.Run("start_time and description set", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceSuspension().Schema, map[string]interface{}{
			"service_id":  1,
			"end_time":    "2099-01-01T00:00:00Z",
			"start_time":  "2098-01-01T00:00:00Z",
			"description": "planned maintenance",
		})

		suspension := suspensionFromResourceData(d)

		if suspension.StartTime == nil || *suspension.StartTime != "2098-01-01T00:00:00Z" {
			t.Errorf("expected start_time to round-trip, got %v", suspension.StartTime)
		}
		if suspension.Description == nil || *suspension.Description != "planned maintenance" {
			t.Errorf("expected description to round-trip, got %v", suspension.Description)
		}
	})
}

func TestSuspensionCreateUpdateDeleteLifecycle(t *testing.T) {
	var stored client.Suspension

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/suspensions":
			json.NewDecoder(r.Body).Decode(&stored)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]int{"id": 1})
		case r.Method == http.MethodGet && r.URL.Path == "/suspensions/1":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(stored)
		case r.Method == http.MethodPut && r.URL.Path == "/suspensions/1":
			json.NewDecoder(r.Body).Decode(&stored)
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodDelete && r.URL.Path == "/suspensions/1":
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	c := client.NewTestClient(server.URL)

	d := schema.TestResourceDataRaw(t, resourceSuspension().Schema, map[string]interface{}{
		"service_id": 1,
		"end_time":   "2099-01-01T00:00:00Z",
	})

	if diags := resourceSuspensionCreate(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error creating suspension: %v", diags)
	}
	if d.Id() != "1" {
		t.Fatalf("expected id \"1\", got %q", d.Id())
	}
	if stored.EndTime != "2099-01-01T00:00:00Z" {
		t.Errorf("expected create request to carry end_time, server stored %+v", stored)
	}

	if diags := resourceSuspensionRead(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error reading suspension: %v", diags)
	}
	if d.Get("service_id").(int) != 1 {
		t.Errorf("expected service_id to round-trip through read, got %v", d.Get("service_id"))
	}

	if err := d.Set("description", "planned maintenance"); err != nil {
		t.Fatalf("unexpected error setting description: %v", err)
	}
	if diags := resourceSuspensionUpdate(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error updating suspension: %v", diags)
	}
	if stored.Description == nil || *stored.Description != "planned maintenance" {
		t.Errorf("expected update request to carry the new description, server stored %+v", stored.Description)
	}

	if diags := resourceSuspensionDelete(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error deleting suspension: %v", diags)
	}
}
