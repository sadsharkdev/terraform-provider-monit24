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

func TestNotificationAddressCreateUpdateDeleteLifecycle(t *testing.T) {
	var stored client.NotificationAddress

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/notification_addresses":
			json.NewDecoder(r.Body).Decode(&stored)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(client.CreateNotificationAddressResponse{ID: 1})
		case r.Method == http.MethodGet && r.URL.Path == "/notification_addresses/1":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(stored)
		case r.Method == http.MethodPut && r.URL.Path == "/notification_addresses/1":
			json.NewDecoder(r.Body).Decode(&stored)
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodDelete && r.URL.Path == "/notification_addresses/1":
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	c := client.NewTestClient(server.URL)

	d := schema.TestResourceDataRaw(t, resourceNotificationAddress().Schema, map[string]interface{}{
		"address":                 "a@example.com",
		"notification_channel_id": "email",
	})

	if diags := resourceNotificationAddressCreate(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error creating notification_address: %v", diags)
	}
	if d.Id() != "1" {
		t.Fatalf("expected id \"1\", got %q", d.Id())
	}
	if stored.Address != "a@example.com" {
		t.Errorf("expected create request to carry address, server stored %+v", stored)
	}

	if diags := resourceNotificationAddressRead(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error reading notification_address: %v", diags)
	}
	if d.Get("address").(string) != "a@example.com" {
		t.Errorf("expected address to round-trip through read, got %q", d.Get("address"))
	}

	if err := d.Set("description", "updated"); err != nil {
		t.Fatalf("unexpected error setting description: %v", err)
	}
	if diags := resourceNotificationAddressUpdate(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error updating notification_address: %v", diags)
	}
	if stored.Description != "updated" {
		t.Errorf("expected update request to carry the new description, server stored %+v", stored.Description)
	}

	if diags := resourceNotificationAddressDelete(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error deleting notification_address: %v", diags)
	}
}
