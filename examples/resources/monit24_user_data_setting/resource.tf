resource "monit24_subaccount" "client_a" {
  name     = "Client A"
  username = "client-a"

  user_data {
    email_address = "client-a@example.com"
  }
}

resource "monit24_user_data_setting" "dashboard_theme" {
  account_id = monit24_subaccount.client_a.id
  key        = "dashboard_theme"
  value      = "dark"
}
