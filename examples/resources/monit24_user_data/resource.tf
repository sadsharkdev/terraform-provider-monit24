resource "monit24_subaccount" "client_a" {
  name     = "Client A"
  username = "client-a"

  user_data {
    email_address = "client-a@example.com"
  }
}

# Managed separately from the subaccount so it can be updated after
# creation (the subaccount's own user_data block is create-only).
resource "monit24_user_data" "client_a" {
  account_id     = monit24_subaccount.client_a.id
  email_address  = "billing@example.com"
  contact_person = "Jane Doe"
  phone_number   = "+48123456789"
}
