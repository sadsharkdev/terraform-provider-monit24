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
