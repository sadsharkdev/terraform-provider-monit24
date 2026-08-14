resource "monit24_periodic_report_address" "daily" {
  address          = "reports@example.com"
  report_frequency = "daily"
  group_id         = monit24_group.my_group.id
}

resource "monit24_periodic_report_address" "weekly" {
  address          = "weekly-reports@example.com"
  report_frequency = "weekly"
  group_id         = monit24_group.my_group.id
}
