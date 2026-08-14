# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A Terraform provider for the Monit24.pl monitoring service, built on `terraform-plugin-sdk/v2`. It wraps the Monit24 REST API (`https://api.monit24.pl/v3`) and exposes it as Terraform resources.

## Commands

```bash
make build         # go build -o terraform-provider-monit24
make install        # build + install into ~/.terraform.d/plugins/<host>/<ns>/<name>/<version>/<os_arch>
make test           # unit tests: go test ./... (excludes vendor)
make testacc        # acceptance tests: TF_ACC=1 go test ./... -timeout 120m (hits the real API)
make docs           # regenerate ./docs from ./examples and ./templates via tfplugindocs
```

Run a single test:
```bash
go test ./monit24/... -run TestAccService -v
TF_ACC=1 go test ./monit24/... -run TestAccService -v -timeout 120m   # acceptance variant
```

Acceptance tests (`TestAcc*`) require real credentials and talk to the live API — they need `MONIT24_USER`/`MONIT24_PASSWORD` (or `MONIT24_TOKEN`, see Authentication below) set (checked by `preCheck` in `monit24/provider_test.go`). Plain `go test` without `TF_ACC=1` skips them.

## Authentication

The provider supports two mutually exclusive auth modes, resolved in `providerConfigure` (`monit24/provider.go`):

1. **`token`** (env `MONIT24_TOKEN`) — sent as `Authorization: Bearer <token>` on every request (`client.NewTokenClient`). This is the **only** way to authenticate an account that has 2FA enabled — Basic Auth is rejected for such accounts with an ambiguous "incorrect credentials or 2FA configured" error. Create the token from the Monit24 account UI (or `POST /sessions` with username+password, which returns a `token`).
2. **`user`+`password`** (env `MONIT24_USER`/`MONIT24_PASSWORD`) — sent as `Authorization: Basic <base64(user:password)>` (`client.NewBasicAuthClient`). Only works for accounts without 2FA.

If `token` is set, it takes priority and `user`/`password` are ignored entirely. `client/client.go`'s `authorizationHeaderValue(basicAuth, token string) string` picks the header; `client_test.go` unit-tests that selection logic directly (no live API needed).

## Architecture

Two-layer structure, consistently applied per resource type:

- **`client/`** — thin, dependency-free HTTP client for the Monit24 REST API. `client/client.go` holds the shared `Client` struct (`get`/`post`/`put`/`delete` helpers, `ResourceNotFound` error type, `OwnerID()`, and the Basic/Bearer auth selection). Each resource has its own file (`client/service.go`, `client/group.go`, `client/notification_address.go`, `client/account.go`, ...) with a data struct (JSON tags matching the API) and `Create*`/`Read*`/`Update*`/`Delete*` methods.
- **`monit24/`** — the Terraform SDK provider and resources. `monit24/provider.go` defines the provider schema (`user`/`password`) and registers resources in `ResourcesMap`. Each `monit24/resource_*.go` defines the Terraform schema and `CreateContext`/`ReadContext`/`UpdateContext`/`DeleteContext` functions that translate between `*schema.ResourceData` and the corresponding `client` struct.

Key conventions to follow when touching a resource:
- IDs are API-assigned integers, stored in Terraform state as strings via `strconv.Itoa`/`strconv.Atoi`. Where there's no single numeric ID (e.g. `monit24_group_share`), a composite `"<a>:<b>"` string ID is used instead — see `groupShareID`/`parseGroupShareID` in `monit24/resource_group_share.go` for the pattern.
- On `ReadContext`, when the API returns 404 (`client.ResourceNotFound`), call `d.SetId("")` and return `nil` (not an error) so Terraform drops it from state.
- All resources support import via `schema.ImportStatePassthroughContext`.
- Optional/nullable API fields are modeled as pointers (`*string`, `*bool`, `*[]int`, etc.) in the `client` structs, only set in `Read*` when non-nil.
- `resourceServiceCreate` delegates to `resourceServiceUpdate` after creation to populate all fields in one pass — follow this pattern for new resources with many optional fields, *unless* Update has side effects beyond a plain PUT (see `monit24_subaccount`, which deliberately does not delegate Create to Update because Update also conditionally calls the `change_password` action).
- `client.Client.OwnerID()` resolves to the parent account ID when the authenticated user is a sub-account; resources pass this as `owner_id` on create/update.
- `extended_settings` (service resource) is a free-form string map; values are coerced to int/bool/string on write (`newServiceFromResourceData`) and merged against currently-defined keys on read (`mergeMaps`) since the API can return additional settings the config doesn't declare.
- A resource whose underlying API record can't be independently deleted (it's implicitly tied to a parent's lifecycle, e.g. `monit24_user_data` — there's no `DELETE /user_data/{id}`) implements `DeleteContext` as a state-only no-op (`d.SetId(""); return nil`), not an API call.

## Resources

Nine resources are registered in `monit24/provider.go`, tracking Monit24 API v3.51 (see `docs-internal/api-inventory.md` for the full endpoint inventory this was audited against, and `docs/superpowers/specs/2026-08-14-full-api-v3.51-coverage-design.md` for the design driving ongoing work):

- **`monit24_group`** (`resource_group.go` / `client/group.go`) — a container that other resources attach to via `group_id`. Fields: `name` (required), `periodic_daily_reports`/`periodic_weekly_reports`/`periodic_monthly_reports` (optional bool, default `true`), `archived_services_in_periodic_reports` (optional bool, default `true`), `assigned_sensor_ids` (optional set of `{category, sensor_ids}` blocks), `is_default` (computed). Does **not** expose the API's deprecated `sensor_ids` field on `group`.
- **`monit24_group_share`** (`resource_group_share.go` / `client/group_share.go`) — shares a group with another account. Composite ID `"<group_id>:<account_id>"` (no server-assigned numeric ID; `PUT /groups/{group_id}/shares/{account_id}` both creates and updates). Fields: `group_id`/`account_id` (required, `ForceNew`), five `can_*` permission bools (optional, default `false`).
- **`monit24_notification_address`** (`resource_notification_address.go` / `client/notification_address.go`) — a channel-specific address notifications are sent to. Fields: `address` (required), `notification_channel_id` (required), `group_id` (optional/computed), `description` (optional, default `""`).
- **`monit24_periodic_report_address`** (`resource_periodic_report_address.go` / `client/periodic_report_address.go`) — sibling of `notification_address` for periodic report delivery. Fields: `address` (required, email), `report_frequency` (required, `daily`/`weekly`/`monthly`), `group_id` (optional/computed).
- **`monit24_service`** (`resource_service.go` / `client/service.go`) — a monitored endpoint/check. Fields: `type_id`/`name`/`address` (required), `group_id` (optional/computed), `interval` (optional, default `600`), `description` (optional), `is_active` (optional, default `true`), `sensor_ids` (optional set of ints), `step_names` (optional ordered list of strings), `notification_channel_ids`/`notification_condition_ids` (optional/computed sets of strings), `notification_mode_id`/`recovery_notification_mode_id` (optional, default `"default"`), `is_archived` (optional, default `false`), `extended_settings` (optional/computed free-form string map). Does not expose `silent_hours`/`suspension_hours` — deprecated in favor of `weekly_suspension`. **`is_archived` write semantics are unverified against the live API** — check via `make testacc` before relying on it.
- **`monit24_subaccount`** (`resource_subaccount.go` / `client/account.go`) — creates a dependent/linked account via `POST /accounts/subaccount`. Fields: `name`/`username` (required), `package_id` (optional/computed), `is_read_only`/`disable_legacy_notifications` (optional bool, default `false`), `language_id` (optional, default `"pl"`), `time_zone_id` (optional, default `"europe_warsaw"`), `is_activated`/`is_blocked`/`parent_account_id` (computed), `subaccount_block`/`subaccount_edit`/`is_2fa_setup_required` (optional bool, default `false`, `ForceNew`), `user_data` (required, `ForceNew`, single nested block — creation-only), `password` (optional, sensitive, never read back, changes routed through `change_password` on Update so it doesn't force recreation), `set_password_url` (optional, `ForceNew`). **`parent_account_id`'s "inferred by the API" assumption and `DELETE /accounts/{id}` vs. `close` for removal are both unverified against the live API** — the acceptance test asserts a specific value/behavior so the next real `make testacc` run conclusively proves or disproves them.
- **`monit24_user_data`** (`resource_user_data.go` / `client/user_data.go`) — manages the contact/billing details tied to an account (`PUT /user_data/{id}`) independently of `monit24_subaccount`'s creation-time `user_data` block. Fields: `account_id` (required, `ForceNew` — this *is* the resource's identity, same numeric ID as the account), `email_address` (required), `address`/`contact_person`/`phone_number`/`tax_identification_number` (optional), `ip_whitelist` (optional list of strings), `ip_whitelist_enabled` (optional bool, default `false`). `Delete` is a state-only no-op — there's no `DELETE /user_data/{id}`.
- **`monit24_suspension`** (`resource_suspension.go` / `client/suspension.go`) — a one-off planned maintenance window for a service. Fields: `service_id` (required), `end_time` (required, ISO datetime), `start_time` (optional/computed), `only_notifications` (optional bool, default `false`), `description` (optional).
- **`monit24_weekly_suspension`** (`resource_weekly_suspension.go` / `client/weekly_suspension.go`) — a recurring maintenance window. `start_minute`/`end_minute` are required single-block nested objects (`day_of_week` 1-7, `hour` 0-23, `minute` 0-59) via `minuteOfWeekSchema()`. Also: `service_id` (required), `only_notifications` (optional bool, default `false`), `description` (optional).

### Known gaps / deliberately out of scope

A second "conditional alerting" subsystem exists in the API (`contacts`, `contact_groups`, `contact_addresses`, `events`, `escalations`, plus their suspension variants) — additive to, not a replacement for, `notification_address`; not yet modeled, see the design spec for the planned phased rollout. `report_templates` and `templates` (custom notification templates), `POST /accounts` (adding another login to the same account, distinct from `monit24_subaccount`'s dependent-account model), and the standalone `/user_data/{id}/settings/{key}` endpoint are also not yet modeled — see the design spec's account-family fixes section for the planned additions. `reports` and `corrections` are intentionally never going to be resources (generate-once/point-in-time, no convergeable state). `sessions` is intentionally never going to be a resource (`POST /sessions` needs the same username+password the provider's Basic Auth already uses — no benefit to modeling it). `admin/*` is out of scope (Monit24-staff-only API surface). `Provider().DataSourcesMap` is still empty — the design spec's Phase 1 covers dictionary data sources, Phase 6 covers per-resource companion data sources.
