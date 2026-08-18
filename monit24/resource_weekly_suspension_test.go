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

func TestWeeklySuspensionCreateUpdateDeleteLifecycle(t *testing.T) {
	var stored client.WeeklySuspension

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/weekly_suspensions":
			json.NewDecoder(r.Body).Decode(&stored)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]int{"id": 1})
		case r.Method == http.MethodGet && r.URL.Path == "/weekly_suspensions/1":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(stored)
		case r.Method == http.MethodPut && r.URL.Path == "/weekly_suspensions/1":
			json.NewDecoder(r.Body).Decode(&stored)
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodDelete && r.URL.Path == "/weekly_suspensions/1":
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	c := client.NewTestClient(server.URL)

	d := schema.TestResourceDataRaw(t, resourceWeeklySuspension().Schema, map[string]interface{}{
		"service_id": 1,
		"start_minute": []interface{}{
			map[string]interface{}{"day_of_week": 6, "hour": 22, "minute": 0},
		},
		"end_minute": []interface{}{
			map[string]interface{}{"day_of_week": 7, "hour": 2, "minute": 0},
		},
	})

	if diags := resourceWeeklySuspensionCreate(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error creating weekly_suspension: %v", diags)
	}
	if d.Id() != "1" {
		t.Fatalf("expected id \"1\", got %q", d.Id())
	}
	if stored.StartMinute.DayOfWeek != 6 {
		t.Errorf("expected create request to carry start_minute, server stored %+v", stored)
	}

	if diags := resourceWeeklySuspensionRead(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error reading weekly_suspension: %v", diags)
	}
	if d.Get("service_id").(int) != 1 {
		t.Errorf("expected service_id to round-trip through read, got %v", d.Get("service_id"))
	}

	if err := d.Set("description", "weekend quiet hours"); err != nil {
		t.Fatalf("unexpected error setting description: %v", err)
	}
	if diags := resourceWeeklySuspensionUpdate(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error updating weekly_suspension: %v", diags)
	}
	if stored.Description == nil || *stored.Description != "weekend quiet hours" {
		t.Errorf("expected update request to carry the new description, server stored %+v", stored.Description)
	}

	if diags := resourceWeeklySuspensionDelete(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error deleting weekly_suspension: %v", diags)
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
