resource "monit24_group" "example" {
  name                                   = "An example group"
  periodic_daily_reports                 = true
  periodic_weekly_reports                = true
  periodic_monthly_reports               = true
  archived_services_in_periodic_reports  = true

  assigned_sensor_ids {
    category   = "default"
    sensor_ids = [1, 2]
  }
}

