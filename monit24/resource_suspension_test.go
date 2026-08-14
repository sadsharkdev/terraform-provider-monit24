package monit24

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSuspension(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSuspensionConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_suspension.test", "end_time", "2099-01-01T00:00:00Z"),
					resource.TestCheckResourceAttr("monit24_suspension.test", "only_notifications", "false"),
				),
			},
			{
				Config: testAccSuspensionConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_suspension.test", "end_time", "2099-06-01T00:00:00Z"),
					resource.TestCheckResourceAttr("monit24_suspension.test", "only_notifications", "true"),
					resource.TestCheckResourceAttr("monit24_suspension.test", "description", "planned maintenance"),
				),
			},
		},
	})
}

const testAccSuspensionConfig = `
resource "monit24_suspension" "test" {
  service_id = monit24_service.test.id
  end_time   = "2099-01-01T00:00:00Z"
}

resource "monit24_service" "test" {
  type_id     = "https"
  name        = "https example.com"
  address     = "example.com"
  group_id    = monit24_group.test.id
}

resource "monit24_group" "test" {
  name = "test group"
}
`

const testAccSuspensionConfigUpdated = `
resource "monit24_suspension" "test" {
  service_id          = monit24_service.test.id
  end_time             = "2099-06-01T00:00:00Z"
  only_notifications   = true
  description          = "planned maintenance"
}

resource "monit24_service" "test" {
  type_id     = "https"
  name        = "https example.com"
  address     = "example.com"
  group_id    = monit24_group.test.id
}

resource "monit24_group" "test" {
  name = "test group"
}
`
