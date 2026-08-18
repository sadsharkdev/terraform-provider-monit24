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

func TestAccGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_group.test", "name", "test group"),
					resource.TestCheckResourceAttr("monit24_group.test", "periodic_daily_reports", "true"),
					resource.TestCheckResourceAttr("monit24_group.test", "periodic_weekly_reports", "true"),
					resource.TestCheckResourceAttr("monit24_group.test", "periodic_monthly_reports", "true"),
					resource.TestCheckResourceAttr("monit24_group.test", "archived_services_in_periodic_reports", "true"),
				),
			},
			{
				Config: testAccGroupConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_group.test", "name", "test group"),
					resource.TestCheckResourceAttr("monit24_group.test", "periodic_daily_reports", "false"),
					resource.TestCheckResourceAttr("monit24_group.test", "periodic_weekly_reports", "false"),
					resource.TestCheckResourceAttr("monit24_group.test", "periodic_monthly_reports", "false"),
					resource.TestCheckResourceAttr("monit24_group.test", "archived_services_in_periodic_reports", "false"),
				),
			},
		},
	})
}

const testAccGroupConfig = `
resource "monit24_group" "test" {
  name = "test group"
}
`

const testAccGroupConfigUpdated = `
resource "monit24_group" "test" {
  name                                   = "test group"
  periodic_daily_reports                 = false
  periodic_weekly_reports                = false
  periodic_monthly_reports               = false
  archived_services_in_periodic_reports  = false
}
`

func TestGroupFromResourceDataAssignedSensorIDs(t *testing.T) {
	t.Run("duplicate category is rejected", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceGroup().Schema, map[string]interface{}{
			"name": "test group",
			"assigned_sensor_ids": []interface{}{
				map[string]interface{}{
					"category":   "default",
					"sensor_ids": []interface{}{1, 2},
				},
				map[string]interface{}{
					"category":   "default",
					"sensor_ids": []interface{}{3, 4},
				},
			},
		})

		_, err := groupFromResourceData(d, client.Client{})
		if err == nil {
			t.Fatal("expected an error for duplicate assigned_sensor_ids category, got nil")
		}
	})

	t.Run("distinct categories are all kept", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceGroup().Schema, map[string]interface{}{
			"name": "test group",
			"assigned_sensor_ids": []interface{}{
				map[string]interface{}{
					"category":   "default",
					"sensor_ids": []interface{}{1, 2},
				},
				map[string]interface{}{
					"category":   "backup",
					"sensor_ids": []interface{}{3},
				},
			},
		})

		group, err := groupFromResourceData(d, client.Client{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if group.AssignedSensorIDs == nil {
			t.Fatal("expected AssignedSensorIDs to be set")
		}

		assigned := *group.AssignedSensorIDs
		if len(assigned) != 2 {
			t.Fatalf("expected 2 categories, got %d: %v", len(assigned), assigned)
		}
	})

	t.Run("field not configured is not sent", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceGroup().Schema, map[string]interface{}{
			"name": "test group",
		})

		group, err := groupFromResourceData(d, client.Client{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if group.AssignedSensorIDs != nil {
			t.Fatalf("expected AssignedSensorIDs to stay nil when never configured, got %v", *group.AssignedSensorIDs)
		}
	})
}

func TestGroupCreateUpdateDeleteLifecycle(t *testing.T) {
	var stored client.Group

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/groups":
			if err := json.NewDecoder(r.Body).Decode(&stored); err != nil {
				t.Fatalf("failed to decode POST body: %v", err)
			}
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(client.CreateGroupResponse{ID: 1})
		case r.Method == http.MethodGet && r.URL.Path == "/groups/1":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(stored)
		case r.Method == http.MethodPut && r.URL.Path == "/groups/1":
			if err := json.NewDecoder(r.Body).Decode(&stored); err != nil {
				t.Fatalf("failed to decode PUT body: %v", err)
			}
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodDelete && r.URL.Path == "/groups/1":
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	c := client.NewTestClient(server.URL)

	d := schema.TestResourceDataRaw(t, resourceGroup().Schema, map[string]interface{}{
		"name": "test group",
	})

	if diags := resourceGroupCreate(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error creating group: %v", diags)
	}
	if d.Id() != "1" {
		t.Fatalf("expected id \"1\", got %q", d.Id())
	}
	if stored.Name != "test group" {
		t.Errorf("expected create request to carry name, server stored %+v", stored)
	}
	if stored.AssignedSensorIDs != nil {
		t.Errorf("expected assigned_sensor_ids to be omitted on create when never configured, server stored %v", *stored.AssignedSensorIDs)
	}

	// Configure assigned_sensor_ids for the first time and update — this is
	// the HasChange("assigned_sensor_ids")-gated path in groupFromResourceData.
	// d.Set doesn't affect HasChange (it reads the diff, not the setWriter
	// overlay), so a fresh ResourceData built from raw config is used here to
	// get a genuine "configured" diff, with the id carried over manually.
	d = schema.TestResourceDataRaw(t, resourceGroup().Schema, map[string]interface{}{
		"name": "test group",
		"assigned_sensor_ids": []interface{}{
			map[string]interface{}{
				"category":   "default",
				"sensor_ids": []interface{}{1, 2},
			},
		},
	})
	d.SetId("1")

	if diags := resourceGroupUpdate(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error updating group: %v", diags)
	}
	if stored.AssignedSensorIDs == nil {
		t.Fatal("expected update request to carry assigned_sensor_ids after configuring it")
	}
	if ids := (*stored.AssignedSensorIDs)["default"]; len(ids) != 2 {
		t.Errorf("expected assigned_sensor_ids[default] to have 2 ids, got %v", ids)
	}

	if diags := resourceGroupDelete(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error deleting group: %v", diags)
	}
}
