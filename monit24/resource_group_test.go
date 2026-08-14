package monit24

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_group.test", "name", "test group"),
					resource.TestCheckResourceAttr("monit24_group.test", "periodic_daily_reports", "true"),
					resource.TestCheckResourceAttr("monit24_group.test", "periodic_weekly_reports", "true"),
					resource.TestCheckResourceAttr("monit24_group.test", "periodic_monthly_reports", "true"),
					resource.TestCheckResourceAttr("monit24_group.test", "archived_services_in_periodic_reports", "true"),
				),
			},
			{
				Config: testAccGroupConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_group.test", "name", "test group"),
					resource.TestCheckResourceAttr("monit24_group.test", "periodic_daily_reports", "false"),
					resource.TestCheckResourceAttr("monit24_group.test", "periodic_weekly_reports", "false"),
					resource.TestCheckResourceAttr("monit24_group.test", "periodic_monthly_reports", "false"),
					resource.TestCheckResourceAttr("monit24_group.test", "archived_services_in_periodic_reports", "false"),
				),
			},
		},
	})
}

const testAccGroupConfig = `
resource "monit24_group" "test" {
  name = "test group"
}
`

const testAccGroupConfigUpdated = `
resource "monit24_group" "test" {
  name                                   = "test group"
  periodic_daily_reports                 = false
  periodic_weekly_reports                = false
  periodic_monthly_reports               = false
  archived_services_in_periodic_reports  = false
}
`
