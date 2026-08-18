package monit24

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccWeeklySuspension(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWeeklySuspensionConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_weekly_suspension.test", "start_minute.0.day_of_week", "6"),
					resource.TestCheckResourceAttr("monit24_weekly_suspension.test", "start_minute.0.hour", "22"),
					resource.TestCheckResourceAttr("monit24_weekly_suspension.test", "end_minute.0.day_of_week", "7"),
					resource.TestCheckResourceAttr("monit24_weekly_suspension.test", "end_minute.0.hour", "2"),
					resource.TestCheckResourceAttr("monit24_weekly_suspension.test", "only_notifications", "false"),
				),
			},
			{
				Config: testAccWeeklySuspensionConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_weekly_suspension.test", "only_notifications", "true"),
					resource.TestCheckResourceAttr("monit24_weekly_suspension.test", "description", "weekend quiet hours"),
				),
			},
		},
	})
}

const testAccWeeklySuspensionConfig = `
resource "monit24_weekly_suspension" "test" {
  service_id = monit24_service.test.id

  start_minute {
    day_of_week = 6
    hour        = 22
    minute      = 0
  }

  end_minute {
    day_of_week = 7
    hour        = 2
    minute      = 0
  }
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

func TestWeeklySuspensionFromResourceData(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceWeeklySuspension().Schema, map[string]interface{}{
		"service_id": 1,
		"start_minute": []interface{}{
			map[string]interface{}{"day_of_week": 6, "hour": 22, "minute": 0},
		},
		"end_minute": []interface{}{
			map[string]interface{}{"day_of_week": 7, "hour": 2, "minute": 0},
		},
	})

	suspension := weeklySuspensionFromResourceData(d)

	if suspension.ServiceID != 1 {
		t.Errorf("expected service_id=1, got %d", suspension.ServiceID)
	}
	if suspension.StartMinute.DayOfWeek != 6 || suspension.StartMinute.Hour != 22 {
		t.Errorf("expected start_minute to round-trip, got %+v", suspension.StartMinute)
	}
	if suspension.EndMinute.DayOfWeek != 7 || suspension.EndMinute.Hour != 2 {
		t.Errorf("expected end_minute to round-trip, got %+v", suspension.EndMinute)
	}
	if suspension.Description != nil {
		t.Errorf("expected description to stay nil when not configured, got %v", *suspension.Description)
	}
}

func TestMinuteOfWeekResourceDataRoundTrip(t *testing.T) {
	raw := []interface{}{
		map[string]interface{}{"day_of_week": 3, "hour": 14, "minute": 30},
	}

	m := minuteOfWeekFromResourceData(raw)
	if m.DayOfWeek != 3 || m.Hour != 14 || m.Minute != 30 {
		t.Fatalf("expected {3 14 30}, got %+v", m)
	}

	back := minuteOfWeekToResourceData(m)
	got := back[0].(map[string]interface{})
	if got["day_of_week"] != 3 || got["hour"] != 14 || got["minute"] != 30 {
		t.Fatalf("expected round-trip to preserve values, got %+v", got)
	}
}

const testAccWeeklySuspensionConfigUpdated = `
resource "monit24_weekly_suspension" "test" {
  service_id           = monit24_service.test.id
  only_notifications    = true
  description           = "weekend quiet hours"

  start_minute {
    day_of_week = 6
    hour        = 22
    minute      = 0
  }

  end_minute {
    day_of_week = 7
    hour        = 2
    minute      = 0
  }
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
