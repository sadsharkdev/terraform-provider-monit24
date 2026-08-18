package monit24

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
