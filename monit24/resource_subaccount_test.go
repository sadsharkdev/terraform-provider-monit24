package monit24

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/monit24/terraform-provider-monit24/client"
)

func TestAccSubaccount(t *testing.T) {
	var subaccountID int

	username := fmt.Sprintf("tf-acc-test-%d", time.Now().UnixNano())

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckSubaccountDestroyed(&subaccountID),
		Steps: []resource.TestStep{
			{
				Config: testAccSubaccountConfig(username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_subaccount.test", "name", "TF acceptance test subaccount"),
					resource.TestCheckResourceAttr("monit24_subaccount.test", "username", username),
					resource.TestCheckResourceAttr("monit24_subaccount.test", "user_data.0.email_address", "tf-acceptance-tests@example.com"),
					func(state *terraform.State) error {
						c, err := client.NewBasicAuthClient(context.Background(), os.Getenv("MONIT24_USER"), os.Getenv("MONIT24_PASSWORD"))
						if err != nil {
							return err
						}

						sub := state.RootModule().Resources["monit24_subaccount.test"]
						gotParentID := sub.Primary.Attributes["parent_account_id"]
						wantParentID := strconv.Itoa(c.OwnerID())

						if gotParentID != wantParentID {
							return fmt.Errorf("expected parent_account_id to equal the caller's own account id (%s) since it's assumed to be inferred by POST /accounts/subaccount, got %s — if this fails, parent_account_id is NOT auto-inferred; make it an Optional (not Computed-only) field in resource_subaccount.go and set it explicitly from OwnerID()", wantParentID, gotParentID)
						}

						return nil
					},
					func(state *terraform.State) error {
						sub := state.RootModule().Resources["monit24_subaccount.test"]
						id, err := strconv.Atoi(sub.Primary.ID)
						if err != nil {
							return err
						}

						subaccountID = id

						return nil
					},
				),
			},
			{
				// Setting a password on an existing subaccount must go through
				// change_password, not a full PUT/recreate — the resource id
				// must stay the same as the previous step.
				Config: testAccSubaccountConfigWithPassword(username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_subaccount.test", "name", "TF acceptance test subaccount"),
					func(state *terraform.State) error {
						sub := state.RootModule().Resources["monit24_subaccount.test"]
						id, err := strconv.Atoi(sub.Primary.ID)
						if err != nil {
							return err
						}

						if id != subaccountID {
							return fmt.Errorf("expected subaccount id to stay %d after changing password, got %d (resource was recreated)", subaccountID, id)
						}

						return nil
					},
				),
			},
		},
	})
}

func testAccCheckSubaccountDestroyed(subaccountID *int) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		c, err := client.NewBasicAuthClient(context.Background(), os.Getenv("MONIT24_USER"), os.Getenv("MONIT24_PASSWORD"))
		if err != nil {
			return err
		}

		_, err = c.ReadAccount(context.Background(), *subaccountID)
		if _, ok := err.(client.ResourceNotFound); ok {
			return nil
		}
		if err != nil {
			return err
		}

		return fmt.Errorf("subaccount %d still exists after destroy", *subaccountID)
	}
}

func testAccSubaccountConfig(username string) string {
	return fmt.Sprintf(`
resource "monit24_subaccount" "test" {
  name     = "TF acceptance test subaccount"
  username = %q

  user_data {
    email_address = "tf-acceptance-tests@example.com"
  }
}
`, username)
}

func testAccSubaccountConfigWithPassword(username string) string {
	return fmt.Sprintf(`
resource "monit24_subaccount" "test" {
  name     = "TF acceptance test subaccount"
  username = %q
  password = "correct-horse-battery-staple"

  user_data {
    email_address = "tf-acceptance-tests@example.com"
  }
}
`, username)
}

func TestAccountFromResourceData(t *testing.T) {
	t.Run("required fields and package_id unset stays nil", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceSubaccount().Schema, map[string]interface{}{
			"name":     "example",
			"username": "example-user",
		})

		account := accountFromResourceData(d)

		if account.Name != "example" || account.Username != "example-user" {
			t.Errorf("expected name/username to round-trip, got %+v", account)
		}
		if account.PackageID != nil {
			t.Errorf("expected package_id to stay nil when not configured, got %v", *account.PackageID)
		}
		if account.IsReadOnly == nil || *account.IsReadOnly {
			t.Errorf("expected is_read_only default false, got %v", account.IsReadOnly)
		}
		if account.DisableLegacyNotifications == nil || *account.DisableLegacyNotifications {
			t.Errorf("expected disable_legacy_notifications default false, got %v", account.DisableLegacyNotifications)
		}
		if account.LanguageID == nil || *account.LanguageID != "pl" {
			t.Errorf("expected language_id default \"pl\", got %v", account.LanguageID)
		}
		if account.TimeZoneID == nil || *account.TimeZoneID != "europe_warsaw" {
			t.Errorf("expected time_zone_id default \"europe_warsaw\", got %v", account.TimeZoneID)
		}
	})

	t.Run("package_id set is passed through as a pointer", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceSubaccount().Schema, map[string]interface{}{
			"name":       "example",
			"username":   "example-user",
			"package_id": 5,
		})

		account := accountFromResourceData(d)

		if account.PackageID == nil || *account.PackageID != 5 {
			t.Errorf("expected package_id=5, got %v", account.PackageID)
		}
	})
}

func TestUserDataFromResourceData(t *testing.T) {
	t.Run("required field only", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceSubaccount().Schema, map[string]interface{}{
			"name":     "example",
			"username": "example-user",
			"user_data": []interface{}{
				map[string]interface{}{
					"email_address": "test@example.com",
				},
			},
		})

		userData := userDataFromResourceData(d)

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

	t.Run("optional fields and ip_whitelist", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceSubaccount().Schema, map[string]interface{}{
			"name":     "example",
			"username": "example-user",
			"user_data": []interface{}{
				map[string]interface{}{
					"email_address":             "test@example.com",
					"address":                   "1 Example St",
					"contact_person":            "Jane Doe",
					"phone_number":              "+48123456789",
					"tax_identification_number": "1234567890",
					"ip_whitelist":              []interface{}{"1.2.3.4", "5.6.7.8"},
					"ip_whitelist_enabled":      true,
				},
			},
		})

		userData := userDataFromResourceData(d)

		if userData.Address == nil || *userData.Address != "1 Example St" {
			t.Errorf("expected address to round-trip, got %v", userData.Address)
		}
		if userData.ContactPerson == nil || *userData.ContactPerson != "Jane Doe" {
			t.Errorf("expected contact_person to round-trip, got %v", userData.ContactPerson)
		}
		if userData.PhoneNumber == nil || *userData.PhoneNumber != "+48123456789" {
			t.Errorf("expected phone_number to round-trip, got %v", userData.PhoneNumber)
		}
		if userData.TaxIdentificationNumber == nil || *userData.TaxIdentificationNumber != "1234567890" {
			t.Errorf("expected tax_identification_number to round-trip, got %v", userData.TaxIdentificationNumber)
		}
		if userData.IPWhitelist == nil || len(*userData.IPWhitelist) != 2 || (*userData.IPWhitelist)[0] != "1.2.3.4" {
			t.Errorf("expected ordered ip_whitelist, got %v", userData.IPWhitelist)
		}
		if userData.IPWhitelistEnabled == nil || !*userData.IPWhitelistEnabled {
			t.Errorf("expected ip_whitelist_enabled=true, got %v", userData.IPWhitelistEnabled)
		}
	})
}
