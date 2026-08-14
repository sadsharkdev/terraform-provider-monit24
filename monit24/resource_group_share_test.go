package monit24

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccGroupShare(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			if os.Getenv("MONIT24_SHARE_ACCOUNT_ID") == "" {
				t.Skip("MONIT24_SHARE_ACCOUNT_ID must be set to an account ID the test account can share groups with")
			}
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupShareConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_group_share.test", "can_modify_group", "false"),
					resource.TestCheckResourceAttr("monit24_group_share.test", "can_create_services", "false"),
				),
			},
			{
				Config: testAccGroupShareConfigUpdated(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_group_share.test", "can_modify_group", "true"),
					resource.TestCheckResourceAttr("monit24_group_share.test", "can_create_services", "true"),
				),
			},
		},
	})
}

func testAccGroupShareConfig() string {
	accountID, _ := strconv.Atoi(os.Getenv("MONIT24_SHARE_ACCOUNT_ID"))

	return fmt.Sprintf(`
resource "monit24_group_share" "test" {
  group_id   = monit24_group.test.id
  account_id = %d
}

resource "monit24_group" "test" {
  name = "test group"
}
`, accountID)
}

func testAccGroupShareConfigUpdated() string {
	accountID, _ := strconv.Atoi(os.Getenv("MONIT24_SHARE_ACCOUNT_ID"))

	return fmt.Sprintf(`
resource "monit24_group_share" "test" {
  group_id              = monit24_group.test.id
  account_id            = %d
  can_modify_group      = true
  can_create_services   = true
}

resource "monit24_group" "test" {
  name = "test group"
}
`, accountID)
}
