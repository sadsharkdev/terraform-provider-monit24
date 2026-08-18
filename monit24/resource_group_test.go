package monit24

import (
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
