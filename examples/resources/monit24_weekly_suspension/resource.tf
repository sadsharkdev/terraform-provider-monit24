resource "monit24_weekly_suspension" "weekend_quiet_hours" {
  service_id          = monit24_service.https_service.id
  only_notifications  = true
  description         = "No notifications during weekend maintenance window"

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
