resource "monit24_suspension" "maintenance_window" {
  service_id          = monit24_service.https_service.id
  start_time          = "2026-01-01T22:00:00Z"
  end_time            = "2026-01-02T02:00:00Z"
  only_notifications  = false
  description         = "Planned maintenance window"
}
