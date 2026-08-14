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

Acceptance tests (`TestAcc*`) require real credentials and talk to the live API — they need `MONIT24_USER` and `MONIT24_PASSWORD` set (checked by `preCheck` in `monit24/provider_test.go`). Plain `go test` without `TF_ACC=1` skips them.

## Architecture

Two-layer structure, consistently applied per resource type (group, notification_address, service):

- **`client/`** — thin, dependency-free HTTP client for the Monit24 REST API. `client/client.go` holds the shared `Client` struct (basic-auth HTTP wrapper with `get`/`post`/`put`/`delete` helpers, `ResourceNotFound` error type, and `OwnerID()`). Each resource has its own file (`client/service.go`, `client/group.go`, `client/notification_address.go`) with a data struct (JSON tags matching the API) and `Create*`/`Read*`/`Update*`/`Delete*` methods.
- **`monit24/`** — the Terraform SDK provider and resources. `monit24/provider.go` defines the provider schema (`user`/`password`, env-fallback via `MONIT24_USER`/`MONIT24_PASSWORD`) and registers resources in `ResourcesMap`. Each `monit24/resource_*.go` defines the Terraform schema and `CreateContext`/`ReadContext`/`UpdateContext`/`DeleteContext` functions that translate between `*schema.ResourceData` and the corresponding `client` struct.

Key conventions to follow when touching a resource:
- IDs are API-assigned integers, stored in Terraform state as strings via `strconv.Itoa`/`strconv.Atoi`.
- On `ReadContext`, when the API returns 404 (`client.ResourceNotFound`), call `d.SetId("")` and return `nil` (not an error) so Terraform drops it from state.
- All resources support import via `schema.ImportStatePassthroughContext`.
- Optional/nullable API fields are modeled as pointers (`*string`, `*bool`, `*[]int`, etc.) in the `client` structs, only set in `Read*` when non-nil.
- `resourceServiceCreate` delegates to `resourceServiceUpdate` after creation to populate all fields in one pass — follow this pattern for new resources with many optional fields rather than duplicating field-setting logic.
- `client.Client.OwnerID()` resolves to the parent account ID when the authenticated user is a sub-account; resources pass this as `owner_id` on create/update.
- `extended_settings` (service resource) is a free-form string map; values are coerced to int/bool/string on write (`newServiceFromResourceData`) and merged against currently-defined keys on read (`mergeMaps`) since the API can return additional settings the config doesn't declare.

## Resources

Seven resources are registered in `monit24/provider.go`, tracking Monit24 API v3.51 (see `docs-internal/api-inventory.md` for the full API endpoint inventory this was audited against):

- **`monit24_group`** (`resource_group.go` / `client/group.go`) — a container that other resources attach to via `group_id`. Fields: `name` (required), `periodic_daily_reports`/`periodic_weekly_reports`/`periodic_monthly_reports` (optional bool, default `true`), `archived_services_in_periodic_reports` (optional bool, default `true`), `assigned_sensor_ids` (optional set of `{category, sensor_ids}` blocks, mapping to the API's `map[string][]int`), `is_default` (computed). Deliberately does **not** expose the API's `sensor_ids` field on `group` — it's flagged for future removal in favor of `assigned_sensor_ids`.
- **`monit24_group_share`** (`resource_group_share.go` / `client/group_share.go`) — shares a group with another account. No server-assigned numeric ID; the API path is `/groups/{group_id}/shares/{account_id}` and `PUT` both creates and updates, so this resource uses a composite Terraform ID (`"<group_id>:<account_id>"`, see `groupShareID`/`parseGroupShareID`) instead of the `strconv.Itoa`-on-create pattern the other resources use. Fields: `group_id`/`account_id` (required, `ForceNew`), five `can_*` permission bools (optional, default `false`).
- **`monit24_notification_address`** (`resource_notification_address.go` / `client/notification_address.go`) — a channel-specific address (e.g. email, phone) that notifications are sent to. Fields: `address` (required), `notification_channel_id` (required, e.g. `"email"`), `group_id` (optional/computed), `description` (optional, default `""`).
- **`monit24_periodic_report_address`** (`resource_periodic_report_address.go` / `client/periodic_report_address.go`) — sibling of `notification_address` for periodic report delivery. Fields: `address` (required, email), `report_frequency` (required, `daily`/`weekly`/`monthly`), `group_id` (optional/computed).
- **`monit24_service`** (`resource_service.go` / `client/service.go`) — a monitored endpoint/check; by far the most complex resource. Fields: `type_id` (required, e.g. `"https"`), `name` (required), `address` (required), `group_id` (optional/computed), `interval` (optional, default `600`), `description` (optional), `is_active` (optional, default `true`), `sensor_ids` (optional set of ints), `step_names` (optional ordered list of strings, for scenario/multi-step services — modeled as `TypeList`, not `TypeSet`, because order matters), `notification_channel_ids` (optional/computed set of strings), `notification_condition_ids` (optional/computed set of strings), `notification_mode_id` (optional, default `"default"`), `recovery_notification_mode_id` (optional, default `"default"`), `extended_settings` (optional/computed free-form string map, e.g. custom sensor params like `http_method`). Deliberately does **not** expose `is_archived` — it was previously removed as unsupported and the API's current write behavior for it hasn't been re-verified against a live account; archiving/restoring should go through the (not-yet-modeled) `/services/{id}/archive` and `/services/{id}/restore` actions. Also does not expose `silent_hours`/`suspension_hours` — deprecated in favor of `weekly_suspension`.
- **`monit24_suspension`** (`resource_suspension.go` / `client/suspension.go`) — a one-off planned maintenance window for a service, replacing the now-deprecated `scheduled_suspensions` API. Fields: `service_id` (required), `end_time` (required, ISO datetime), `start_time` (optional/computed), `only_notifications` (optional bool, default `false`), `description` (optional).
- **`monit24_weekly_suspension`** (`resource_weekly_suspension.go` / `client/weekly_suspension.go`) — a recurring maintenance window, the documented replacement for `monit24_service`'s removed `silent_hours`/`suspension_hours` fields. `start_minute`/`end_minute` are required single-block nested objects (`day_of_week` 1-7, `hour` 0-23, `minute` 0-59 — mirroring the API's `minute_of_week` object, not a raw integer), via the shared `minuteOfWeekSchema()` helper. Fields also include `service_id` (required), `only_notifications` (optional bool, default `false`), `description` (optional).

All resources with a `group_id` field follow the same shape: required fields plus an `owner_id` set server-side from `client.Client.OwnerID()` rather than exposed in the Terraform schema. When adding a new resource, mirror this pair of files (`client/<resource>.go` + `monit24/resource_<resource>.go`), register it in `Provider().ResourcesMap`, add example `.tf` + `import.sh` files under `examples/resources/monit24_<resource>/`, add a `monit24/resource_<resource>_test.go` acceptance test, and run `make docs`. For resources without a server-assigned integer ID (like `monit24_group_share`), use a composite string ID instead of the usual `strconv.Itoa`-after-create pattern.

### Known gaps / deliberately out of scope

Per the `docs-internal/api-inventory.md` audit against Monit24 API v3.51, a second "conditional alerting" subsystem exists (`contacts`, `contact_groups`, `contact_addresses`, `events`, `escalations`, plus their suspension variants) that is **not** modeled by this provider — it's additive to, not a replacement for, `notification_address`. `report_templates`, `templates` (custom notification templates), and `accounts`/subaccounts (CRUD exists but carries sensitive fields like `password`) are also unmodeled. `reports` and `corrections` are intentionally not resources — they're generate-once/point-in-time API objects with no natural "current state" to converge on, a poor fit for Terraform's plan/diff model. `Provider().DataSourcesMap` is still empty; the API's read-only `dictionaries/*` endpoints (service types, notification channels/conditions/modes, time zones, etc.) would be reasonable data source candidates.

## Documentation

`docs/` is generated, not hand-edited — update `examples/` (per-resource `.tf` files) and `templates/*.tmpl`, then run `make docs`.
