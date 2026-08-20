# Monit24.pl Terraform Provider

See full documentation on the [Terraform Registry docs](https://registry.terraform.io/providers/monit24/monit24/latest/docs).

## Usage

```tf
terraform {
  required_providers {
    monit24 = {
      version = ">= 0.1.4"
      source  = "monit24/monit24"
    }
  }
}

provider "monit24" {}

resource "monit24_group" "example" {
  name = "An example group"
}

resource "monit24_notification_address" "email" {
  address                 = "notifications@example.com"
  notification_channel_id = "email"
  group_id                = monit24_group.example.id
  description             = "An example notification email address"
}

resource "monit24_service" "test_https_service" {
  type_id     = "https"
  name        = "https example.com"
  description = "An example https service"
  address     = "example.com"
  group_id    = monit24_group.example.id
}
```

See [examples](./examples) for in-depth usage.

## Authentication

Fill the `user` and `password` directly in the provider (preferably via Terraform's variables).

```tf
provider "monit24" {
    user     = "username"
    password = "topsecret"
}
```

Alternatively, provide `MONIT24_USER` and `MONIT24_PASSWORD` environment variables.

```tf
provider "monit24" {}
```

```bash
export MONIT24_USER=
export MONIT24_PASSWORD=
```

If the account has 2FA enabled, `user`/`password` will not work (the API rejects Basic Auth for 2FA accounts). Use an API token instead — create one from the Monit24 account UI, then:

```tf
provider "monit24" {
    token = "your-api-token"
}
```

```bash
export MONIT24_TOKEN=
```

`token` takes priority over `user`/`password` if both are set.

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md)
