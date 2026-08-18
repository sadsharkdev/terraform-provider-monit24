package monit24

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/monit24/terraform-provider-monit24/client"
)

func TestAccNotificationAddress(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNotificationAddressConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_notification_address.test", "address", "notifications@example.com"),
				),
			},
		},
	})
}

const testAccNotificationAddressConfig = `
resource "monit24_notification_address" "test" {
  address                 = "notifications@example.com"
  notification_channel_id = "email"
  group_id                = monit24_group.test.id
  description             = "Email notification"
}

resource "monit24_group" "test" {
  name = "test group"
}
`

func TestNotificationAddressFromResourceData(t *testing.T) {
	t.Run("required fields and description omitted", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceNotificationAddress().Schema, map[string]interface{}{
			"address":                 "a@example.com",
			"notification_channel_id": "email",
			"group_id":                5,
		})

		address := notificationAddressFromResourceData(d, client.Client{})

		if address.Address != "a@example.com" || address.NotificationChannelID != "email" {
			t.Errorf("expected address/notification_channel_id to round-trip, got %+v", address)
		}
		if address.GroupID != 5 {
			t.Errorf("expected group_id=5, got %d", address.GroupID)
		}
		if address.Description != "" {
			t.Errorf("expected description to default to empty string, got %q", address.Description)
		}
	})

	t.Run("description set", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceNotificationAddress().Schema, map[string]interface{}{
			"address":                 "a@example.com",
			"notification_channel_id": "email",
			"description":             "primary contact",
		})

		address := notificationAddressFromResourceData(d, client.Client{})

		if address.Description != "primary contact" {
			t.Errorf("expected description to round-trip, got %q", address.Description)
		}
	})
}
