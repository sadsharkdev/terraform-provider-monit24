package monit24

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
