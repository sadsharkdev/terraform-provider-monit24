package monit24

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUserDataSetting(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSettingConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_user_data_setting.test", "key", "tf_acceptance_test_setting"),
					resource.TestCheckResourceAttr("monit24_user_data_setting.test", "value", "original-value"),
				),
			},
			{
				Config: testAccUserDataSettingConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_user_data_setting.test", "value", "updated-value"),
				),
			},
		},
	})
}

const testAccUserDataSettingConfig = `
resource "monit24_subaccount" "test" {
  name     = "TF acceptance test user_data_setting subaccount"
  username = "tf-acc-test-uds"

  user_data {
    email_address = "tf-acceptance-tests@example.com"
  }
}

resource "monit24_user_data_setting" "test" {
  account_id = monit24_subaccount.test.id
  key        = "tf_acceptance_test_setting"
  value      = "original-value"
}
`

const testAccUserDataSettingConfigUpdated = `
resource "monit24_subaccount" "test" {
  name     = "TF acceptance test user_data_setting subaccount"
  username = "tf-acc-test-uds"

  user_data {
    email_address = "tf-acceptance-tests@example.com"
  }
}

resource "monit24_user_data_setting" "test" {
  account_id = monit24_subaccount.test.id
  key        = "tf_acceptance_test_setting"
  value      = "updated-value"
}
`

func TestUserDataSettingIDRoundTrip(t *testing.T) {
	id := userDataSettingID(123456, "dashboard_theme")
	if id != "123456:dashboard_theme" {
		t.Fatalf("expected %q, got %q", "123456:dashboard_theme", id)
	}

	accountID, key, err := parseUserDataSettingID(id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if accountID != 123456 || key != "dashboard_theme" {
		t.Fatalf("expected (123456, \"dashboard_theme\"), got (%d, %q)", accountID, key)
	}
}

func TestUserDataSettingIDKeyContainingColon(t *testing.T) {
	// SplitN(id, ":", 2) must only split on the first colon, since the key
	// itself may legitimately contain one.
	accountID, key, err := parseUserDataSettingID("123456:namespace:setting")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if accountID != 123456 || key != "namespace:setting" {
		t.Fatalf("expected (123456, \"namespace:setting\"), got (%d, %q)", accountID, key)
	}
}

func TestParseUserDataSettingIDInvalid(t *testing.T) {
	tests := []string{"", "123456", "abc:dashboard_theme"}

	for _, id := range tests {
		if _, _, err := parseUserDataSettingID(id); err == nil {
			t.Errorf("parseUserDataSettingID(%q): expected an error, got nil", id)
		}
	}
}
