package monit24

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
