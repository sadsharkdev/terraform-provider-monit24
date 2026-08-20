package monit24

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccAccountUser(t *testing.T) {
	username := fmt.Sprintf("tf-acc-test-user-%d", time.Now().UnixNano())

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAccountUserConfig(username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_account_user.test", "name", "TF acceptance test account user"),
					resource.TestCheckResourceAttr("monit24_account_user.test", "username", username),
				),
			},
		},
	})
}

func testAccAccountUserConfig(username string) string {
	return fmt.Sprintf(`
resource "monit24_account_user" "test" {
  name     = "TF acceptance test account user"
  username = %q

  user_data {
    email_address = "tf-acceptance-tests@example.com"
  }
}
`, username)
}

func TestAccountUserSchemaHasNoSubaccountOnlyFields(t *testing.T) {
	s := resourceAccountUser().Schema

	for _, field := range []string{"set_password_url", "subaccount_block", "subaccount_edit", "is_2fa_setup_required", "parent_account_id"} {
		if _, ok := s[field]; ok {
			t.Errorf("expected monit24_account_user schema to not define %q (subaccount-only field), but it does", field)
		}
	}
}

func TestAccountUserSharesAccountFromResourceDataMapping(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceAccountUser().Schema, map[string]interface{}{
		"name":     "example",
		"username": "example-user",
		"user_data": []interface{}{
			map[string]interface{}{
				"email_address": "test@example.com",
			},
		},
	})

	account := accountFromResourceData(d)
	if account.Name != "example" || account.Username != "example-user" {
		t.Errorf("expected name/username to round-trip, got %+v", account)
	}

	userData := userDataFromResourceData(d)
	if userData.EmailAddress != "test@example.com" {
		t.Errorf("expected email_address to round-trip, got %q", userData.EmailAddress)
	}
}
