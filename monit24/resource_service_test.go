package monit24

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
