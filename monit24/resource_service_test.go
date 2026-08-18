package monit24

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/monit24/terraform-provider-monit24/client"
)

func TestAccService(t *testing.T) {
	var serviceID int

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceAttrs("monit24_service.test", testServiceAttributesCreated),
					resource.TestCheckTypeSetElemAttr("monit24_service.test", "notification_channel_ids.*", "email"),
					resource.TestCheckTypeSetElemAttr("monit24_service.test", "notification_channel_ids.*", "sms"),
					resource.TestCheckTypeSetElemAttr("monit24_service.test", "notification_condition_ids.*", "failure"),
					resource.TestCheckTypeSetElemAttr("monit24_service.test", "notification_condition_ids.*", "recovery"),
					func(state *terraform.State) error {
						service := state.RootModule().Resources["monit24_service.test"]
						id, err := strconv.Atoi(service.Primary.ID)
						if err != nil {
							return err
						}

						serviceID = id

						return nil
					},
				),
			},
			{
				Config: testAccServiceConfigUpdateDefaults,
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceAttrs("monit24_service.test", testServiceAttributesDefaultsUpdated),
					resource.TestCheckTypeSetElemAttr("monit24_service.test", "notification_channel_ids.*", "sms"),
					resource.TestCheckTypeSetElemAttr("monit24_service.test", "notification_condition_ids.*", "recovery"),
				),
			},
			{
				Config: testAccServiceConfigExtendedSettings,
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceAttrs("monit24_service.test", testServiceAttributesDefaultsUpdated),
					resource.TestCheckTypeSetElemAttr("monit24_service.test", "notification_channel_ids.*", "sms"),
					resource.TestCheckTypeSetElemAttr("monit24_service.test", "notification_condition_ids.*", "recovery"),
					resource.TestCheckResourceAttr("monit24_service.test", "extended_settings.http_method", "POST"),
				),
			},
			{
				// Remove extended settings
				Config: testAccServiceConfigExtendedSettingsDeleted,
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceAttrs("monit24_service.test", testServiceAttributesDefaultsUpdated),
					resource.TestCheckTypeSetElemAttr("monit24_service.test", "notification_channel_ids.*", "sms"),
					resource.TestCheckTypeSetElemAttr("monit24_service.test", "notification_condition_ids.*", "recovery"),
					resource.TestCheckNoResourceAttr("monit24_service.test", "extended_settings.http_method"),
				),
			},
			{
				SkipFunc: func() (bool, error) {
					c, err := client.NewBasicAuthClient(context.Background(), os.Getenv("MONIT24_USER"), os.Getenv("MONIT24_PASSWORD"))
					if err != nil {
						return false, err
					}

					err = c.DeleteService(context.Background(), serviceID)
					if err != nil {
						return false, err
					}

					return false, nil
				},
				// Apply the same settings on remotely deleted resource
				Config: testAccServiceConfigExtendedSettingsDeleted,
				Check: resource.ComposeTestCheckFunc(
					testCheckResourceAttrs("monit24_service.test", testServiceAttributesDefaultsUpdated),
					resource.TestCheckTypeSetElemAttr("monit24_service.test", "notification_channel_ids.*", "sms"),
					resource.TestCheckTypeSetElemAttr("monit24_service.test", "notification_condition_ids.*", "recovery"),
					resource.TestCheckNoResourceAttr("monit24_service.test", "extended_settings.http_method"),
				),
			},
		},
	})
}

var (
	testServiceAttributesCreated = map[string]string{
		"address":                       "example.com",
		"interval":                      "600",
		"is_active":                     "true",
		"is_archived":                   "false",
		"notification_mode_id":          "default",
		"recovery_notification_mode_id": "default",
	}
	testServiceAttributesDefaultsUpdated = map[string]string{
		"address":                       "new.example.com",
		"interval":                      "700",
		"is_active":                     "false",
		"is_archived":                   "false",
		"notification_mode_id":          "off",
		"recovery_notification_mode_id": "after_30_seconds",
	}
)

func testCheckResourceAttrs(name string, attrs map[string]string) resource.TestCheckFunc {
	var fns []resource.TestCheckFunc

	for k, v := range attrs {
		fns = append(fns, resource.TestCheckResourceAttr(name, k, v))
	}

	return resource.ComposeTestCheckFunc(fns...)
}

const testAccServiceConfig = `
resource "monit24_service" "test" {
  type_id                       = "https"
  name                          = "https example.com"
  description                   = "An example HTTPS service"
  address                       = "example.com"
  group_id                      = monit24_group.test.id
}

resource "monit24_group" "test" {
  name = "test group"
}
`

const testAccServiceConfigUpdateDefaults = `
resource "monit24_service" "test" {
  type_id                       = "https"
  name                          = "https example.com"
  description                   = "An example HTTPS service"
  address                       = "new.example.com"
  group_id                      = monit24_group.test.id
  interval                      = 700
  is_active                     = false
  notification_channel_ids      = ["sms"]
  notification_condition_ids    = ["recovery"]
  notification_mode_id          = "off"
  recovery_notification_mode_id = "after_30_seconds"
}

resource "monit24_group" "test" {
  name = "test group"
}
`

const testAccServiceConfigExtendedSettings = `
resource "monit24_service" "test" {
  type_id                       = "https"
  name                          = "https example.com"
  description                   = "An example HTTPS service"
  address                       = "new.example.com"
  group_id                      = monit24_group.test.id
  interval                      = 700
  is_active                     = false
  notification_channel_ids      = ["sms"]
  notification_condition_ids    = ["recovery"]
  notification_mode_id          = "off"
  recovery_notification_mode_id = "after_30_seconds"
  extended_settings = {
    http_method = "POST"
  }
}

resource "monit24_group" "test" {
  name = "test group"
}
`

const testAccServiceConfigExtendedSettingsDeleted = `
resource "monit24_service" "test" {
  type_id                       = "https"
  name                          = "https example.com"
  description                   = "An example HTTPS service"
  address                       = "new.example.com"
  group_id                      = monit24_group.test.id
  interval                      = 700
  is_active                     = false
  notification_channel_ids      = ["sms"]
  notification_condition_ids    = ["recovery"]
  notification_mode_id          = "off"
  recovery_notification_mode_id = "after_30_seconds"
  extended_settings = {
  }
}

resource "monit24_group" "test" {
  name = "test group"
}
`

func TestNewServiceFromResourceData(t *testing.T) {
	t.Run("required and simple optional fields", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceService().Schema, map[string]interface{}{
			"type_id":     "https",
			"name":        "example",
			"address":     "example.com",
			"group_id":    5,
			"interval":    300,
			"description": "an example",
			"is_active":   false,
			"is_archived": true,
		})

		service, err := newServiceFromResourceData(client.Service{}, d)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if service.Name != "example" || service.Address != "example.com" {
			t.Errorf("expected name/address to round-trip, got %+v", service)
		}
		if service.GroupID != 5 || service.Interval != 300 {
			t.Errorf("expected group_id=5 interval=300, got group_id=%d interval=%d", service.GroupID, service.Interval)
		}
		if service.Description == nil || *service.Description != "an example" {
			t.Errorf("expected description=\"an example\", got %v", service.Description)
		}
		if service.IsActive == nil || *service.IsActive {
			t.Errorf("expected is_active=false, got %v", service.IsActive)
		}
		if service.IsArchived == nil || !*service.IsArchived {
			t.Errorf("expected is_archived=true, got %v", service.IsArchived)
		}
	})

	t.Run("sensor_ids and step_names", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceService().Schema, map[string]interface{}{
			"type_id":    "https",
			"name":       "example",
			"address":    "example.com",
			"sensor_ids": []interface{}{1, 2},
			"step_names": []interface{}{"step one", "step two"},
		})

		service, err := newServiceFromResourceData(client.Service{}, d)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if service.SensorIDs == nil || len(*service.SensorIDs) != 2 {
			t.Errorf("expected 2 sensor_ids, got %v", service.SensorIDs)
		}
		if service.StepNames == nil || len(*service.StepNames) != 2 || (*service.StepNames)[0] != "step one" {
			t.Errorf("expected ordered step_names, got %v", service.StepNames)
		}
	})

	t.Run("extended_settings coerces string values to int, bool, or string", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceService().Schema, map[string]interface{}{
			"type_id": "https",
			"name":    "example",
			"address": "example.com",
			"extended_settings": map[string]interface{}{
				"port":        "443",
				"http_method": "POST",
				"is_ssl":      "true",
			},
		})

		service, err := newServiceFromResourceData(client.Service{}, d)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if service.ExtendedSettings == nil {
			t.Fatal("expected ExtendedSettings to be set")
		}
		settings := *service.ExtendedSettings

		if settings["port"] != 443 {
			t.Errorf("expected port to coerce to int 443, got %v (%T)", settings["port"], settings["port"])
		}
		if settings["is_ssl"] != true {
			t.Errorf("expected is_ssl to coerce to bool true, got %v (%T)", settings["is_ssl"], settings["is_ssl"])
		}
		if settings["http_method"] != "POST" {
			t.Errorf("expected http_method to stay a string \"POST\", got %v (%T)", settings["http_method"], settings["http_method"])
		}
	})
}

func TestParseBool(t *testing.T) {
	tests := []struct {
		input   string
		want    bool
		wantErr bool
	}{
		{"true", true, false},
		{"false", false, false},
		{"yes", false, true},
		{"", false, true},
	}

	for _, tt := range tests {
		got, err := parseBool(tt.input)
		if tt.wantErr {
			if err == nil {
				t.Errorf("parseBool(%q): expected an error, got nil", tt.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("parseBool(%q): unexpected error: %v", tt.input, err)
		}
		if got != tt.want {
			t.Errorf("parseBool(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestMergeMaps(t *testing.T) {
	external := map[string]interface{}{
		"http_method": "POST",
		"port":        443,
		"undeclared":  "should be dropped",
	}
	defined := map[string]interface{}{
		"http_method": "GET",
		"port":        "80",
	}

	got := mergeMaps(external, defined)

	if len(got) != 2 {
		t.Fatalf("expected 2 keys (only those present in defined), got %d: %v", len(got), got)
	}
	if got["http_method"] != "POST" {
		t.Errorf("expected http_method=POST (from external, as a string), got %v", got["http_method"])
	}
	if got["port"] != "443" {
		t.Errorf("expected port=\"443\" (int coerced to string via fmt.Sprintf), got %v", got["port"])
	}
	if _, present := got["undeclared"]; present {
		t.Error("expected undeclared (not present in defined) to be dropped")
	}
}
