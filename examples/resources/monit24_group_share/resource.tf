resource "monit24_group_share" "read_only" {
  group_id   = monit24_group.my_group.id
  account_id = 654321
}

resource "monit24_group_share" "editor" {
  group_id             = monit24_group.my_group.id
  account_id           = 654322
  can_modify_group     = true
  can_modify_services  = true
  can_create_services  = true
  can_delete_services  = true
  can_archive_services = true
  can_force_analyses   = true
}
