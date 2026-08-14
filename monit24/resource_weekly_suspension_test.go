package monit24

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccWeeklySuspension(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWeeklySuspensionConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_weekly_suspension.test", "start_minute.0.day_of_week", "6"),
					resource.TestCheckResourceAttr("monit24_weekly_suspension.test", "start_minute.0.hour", "22"),
					resource.TestCheckResourceAttr("monit24_weekly_suspension.test", "end_minute.0.day_of_week", "7"),
					resource.TestCheckResourceAttr("monit24_weekly_suspension.test", "end_minute.0.hour", "2"),
					resource.TestCheckResourceAttr("monit24_weekly_suspension.test", "only_notifications", "false"),
				),
			},
			{
				Config: testAccWeeklySuspensionConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_weekly_suspension.test", "only_notifications", "true"),
					resource.TestCheckResourceAttr("monit24_weekly_suspension.test", "description", "weekend quiet hours"),
				),
			},
		},
	})
}

const testAccWeeklySuspensionConfig = `
resource "monit24_weekly_suspension" "test" {
  service_id = monit24_service.test.id

  start_minute {
    day_of_week = 6
    hour        = 22
    minute      = 0
  }

  end_minute {
    day_of_week = 7
    hour        = 2
    minute      = 0
  }
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

const testAccWeeklySuspensionConfigUpdated = `
resource "monit24_weekly_suspension" "test" {
  service_id           = monit24_service.test.id
  only_notifications    = true
  description           = "weekend quiet hours"

  start_minute {
    day_of_week = 6
    hour        = 22
    minute      = 0
  }

  end_minute {
    day_of_week = 7
    hour        = 2
    minute      = 0
  }
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
