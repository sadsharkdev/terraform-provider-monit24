package monit24

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPeriodicReportAddress(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPeriodicReportAddressConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_periodic_report_address.test", "address", "reports@example.com"),
					resource.TestCheckResourceAttr("monit24_periodic_report_address.test", "report_frequency", "daily"),
				),
			},
			{
				Config: testAccPeriodicReportAddressConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_periodic_report_address.test", "address", "reports@example.com"),
					resource.TestCheckResourceAttr("monit24_periodic_report_address.test", "report_frequency", "weekly"),
				),
			},
		},
	})
}

const testAccPeriodicReportAddressConfig = `
resource "monit24_periodic_report_address" "test" {
  address           = "reports@example.com"
  report_frequency  = "daily"
  group_id          = monit24_group.test.id
}

resource "monit24_group" "test" {
  name = "test group"
}
`

const testAccPeriodicReportAddressConfigUpdated = `
resource "monit24_periodic_report_address" "test" {
  address           = "reports@example.com"
  report_frequency  = "weekly"
  group_id          = monit24_group.test.id
}

resource "monit24_group" "test" {
  name = "test group"
}
`
