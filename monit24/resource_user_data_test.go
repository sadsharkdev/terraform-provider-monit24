package monit24

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccUserData(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_user_data.test", "email_address", "tf-acceptance-tests@example.com"),
					resource.TestCheckResourceAttr("monit24_user_data.test", "contact_person", "Original Contact"),
				),
			},
			{
				Config: testAccUserDataConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_user_data.test", "email_address", "tf-acceptance-tests@example.com"),
					resource.TestCheckResourceAttr("monit24_user_data.test", "contact_person", "Updated Contact"),
				),
			},
		},
	})
}

const testAccUserDataConfig = `
resource "monit24_subaccount" "test" {
  name     = "TF acceptance test user_data subaccount"
  username = "tf-acc-test-user-data"

  user_data {
    email_address = "tf-acceptance-tests@example.com"
  }
}

resource "monit24_user_data" "test" {
  account_id     = monit24_subaccount.test.id
  email_address  = "tf-acceptance-tests@example.com"
  contact_person = "Original Contact"
}
`

const testAccUserDataConfigUpdated = `
resource "monit24_subaccount" "test" {
  name     = "TF acceptance test user_data subaccount"
  username = "tf-acc-test-user-data"

  user_data {
    email_address = "tf-acceptance-tests@example.com"
  }
}

resource "monit24_user_data" "test" {
  account_id     = monit24_subaccount.test.id
  email_address  = "tf-acceptance-tests@example.com"
  contact_person = "Updated Contact"
}
`

func TestFlatUserDataFromResourceData(t *testing.T) {
	t.Run("required field only", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceUserData().Schema, map[string]interface{}{
			"account_id":    1,
			"email_address": "test@example.com",
		})

		userData := flatUserDataFromResourceData(d)

		if userData.EmailAddress != "test@example.com" {
			t.Errorf("expected email_address to round-trip, got %q", userData.EmailAddress)
		}
		if userData.Address != nil {
			t.Errorf("expected address to stay nil when not configured, got %v", *userData.Address)
		}
		if userData.IPWhitelist != nil {
			t.Errorf("expected ip_whitelist to stay nil when not configured, got %v", *userData.IPWhitelist)
		}
		if userData.IPWhitelistEnabled == nil || *userData.IPWhitelistEnabled {
			t.Errorf("expected ip_whitelist_enabled default false, got %v", userData.IPWhitelistEnabled)
		}
	})

	t.Run("optional fields and ip_whitelist set", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceUserData().Schema, map[string]interface{}{
			"account_id":           1,
			"email_address":        "test@example.com",
			"address":              "1 Example St",
			"ip_whitelist":         []interface{}{"1.2.3.4"},
			"ip_whitelist_enabled": true,
		})

		userData := flatUserDataFromResourceData(d)

		if userData.Address == nil || *userData.Address != "1 Example St" {
			t.Errorf("expected address to round-trip, got %v", userData.Address)
		}
		if userData.IPWhitelist == nil || len(*userData.IPWhitelist) != 1 || (*userData.IPWhitelist)[0] != "1.2.3.4" {
			t.Errorf("expected ip_whitelist to round-trip, got %v", userData.IPWhitelist)
		}
		if userData.IPWhitelistEnabled == nil || !*userData.IPWhitelistEnabled {
			t.Errorf("expected ip_whitelist_enabled=true, got %v", userData.IPWhitelistEnabled)
		}
	})
}
