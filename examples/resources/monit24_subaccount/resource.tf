# Password omitted entirely — the subaccount is created without a usable
# password; set one later out-of-band (e.g. via the change_password action)
# or have the account holder use the "forgot password" flow.
resource "monit24_subaccount" "no_password" {
  name     = "Client A"
  username = "client-a"

  user_data {
    email_address = "client-a@example.com"
  }
}

# Password sourced from a Terraform variable — populate it from a local
# tfvars file, TF_VAR_subaccount_password, Vault, or a cloud Secret Manager
# data source. Rotating the value later updates the password in place via
# the API's change_password action rather than recreating the subaccount.
variable "subaccount_password" {
  type      = string
  sensitive = true
}

resource "monit24_subaccount" "with_password" {
  name     = "Client B"
  username = "client-b"
  password = var.subaccount_password

  user_data {
    email_address = "client-b@example.com"
  }
}

# set_password_url instead of password — Monit24 emails the new subaccount
# holder a link to set their own password, so the secret never passes
# through Terraform config or state at all.
resource "monit24_subaccount" "self_service_password" {
  name              = "Client C"
  username          = "client-c"
  set_password_url  = "https://portal.example.com/set-password"

  user_data {
    email_address = "client-c@example.com"
  }
}
