# A second login for the SAME account (distinct from monit24_subaccount,
# which creates a separate, dependent account).
resource "monit24_account_user" "colleague" {
  name     = "Jane's colleague"
  username = "jane-colleague"

  user_data {
    email_address = "colleague@example.com"
  }
}
