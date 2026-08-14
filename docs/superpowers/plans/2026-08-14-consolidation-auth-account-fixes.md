# Branch Consolidation, Token Auth & Account-Family Fixes Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Consolidate the two in-flight resource branches onto `feature/monit24-api-v3.51-resources`, fix the Basic-Auth-only limitation that breaks 2FA-enabled accounts by adding `MONIT24_TOKEN` Bearer auth, and close the account-family gaps identified in the design spec (`is_archived`, `parent_account_id` verification, `monit24_user_data`, `monit24_account_user`, `monit24_user_data_setting`).

**Architecture:** Existing two-layer pattern (`client/*.go` thin HTTP wrappers + `monit24/resource_*.go` Terraform SDK glue). No new architectural pattern — every task either extends `client/client.go`'s auth handling or adds a resource file pair following the established convention.

**Tech Stack:** Go 1.17, `terraform-plugin-sdk/v2`, standard library only in `client/`.

## Global Constraints

- Every task ends with `gofmt -l -w`, `go build ./...`, `go vet ./...`, `go test ./...` all clean before moving on.
- IDs stored as strings via `strconv.Itoa`/`strconv.Atoi` (or composite `"<a>:<b>"` strings where there's no single numeric ID, per `client/group_share.go`'s precedent).
- 404 (`client.ResourceNotFound`) on Read → `d.SetId(""); return nil`, never an error.
- All resources: `Importer: &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext}`.
- Optional/nullable API fields are `*T` pointers in `client` structs, only `d.Set` in `Read` when non-nil.
- Do not push to git or open a PR — commit locally only, on `feature/monit24-api-v3.51-resources`.
- This is the first of several plans (per the design spec's phasing) — it covers only the spec's "Branch & commit strategy," "Authentication fix," and "Account-family fixes" sections. The "Phased resource + data source rollout" section (dictionaries, contacts, events/escalations, reporting, companion data sources) is out of scope here and gets its own plan(s) later.

---

### Task 1: Consolidate branches — cherry-pick `monit24_subaccount` onto this branch

**Files:**
- Modify (conflict resolution): `monit24/provider.go`, `CLAUDE.md`
- No test files (this is a git operation; verification is a full build/vet/test pass)

**Interfaces:**
- Produces: the branch now has both `b157cee`'s resources (group, group_share, notification_address, periodic_report_address, service, suspension, weekly_suspension) and `fb15f07`'s `monit24_subaccount`, all registered in one `Provider().ResourcesMap`.

- [ ] **Step 1: Confirm current branch and commit to cherry-pick**

Run: `git branch --show-current && git log --oneline -1 feature/monit24-subaccount-resource`
Expected: current branch is `feature/monit24-api-v3.51-resources`, and the log shows `fb15f07 Add monit24_subaccount resource for managing linked/dependent accounts`.

- [ ] **Step 2: Cherry-pick it**

Run: `git cherry-pick fb15f07`
Expected: conflicts in `monit24/provider.go` and `CLAUDE.md` (both commits independently modified/created these files from a common ancestor). `client/account.go`, `monit24/resource_subaccount.go`, `monit24/resource_subaccount_test.go`, and `examples/resources/monit24_subaccount/*` apply cleanly (new files, no conflict).

- [ ] **Step 3: Resolve the `monit24/provider.go` conflict**

Replace the conflicted `ResourcesMap` block with:

```go
		ResourcesMap: map[string]*schema.Resource{
			"monit24_group":                   resourceGroup(),
			"monit24_group_share":             resourceGroupShare(),
			"monit24_notification_address":    resourceNotificationAddress(),
			"monit24_periodic_report_address": resourcePeriodicReportAddress(),
			"monit24_service":                 resourceService(),
			"monit24_subaccount":              resourceSubaccount(),
			"monit24_suspension":              resourceSuspension(),
			"monit24_weekly_suspension":       resourceWeeklySuspension(),
		},
```

Leave the rest of the file (provider schema, `providerConfigure`) untouched for now — Task 3 modifies it further.

- [ ] **Step 4: Resolve the `CLAUDE.md` conflict**

Replace the entire file with:

```markdown
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

- **`client/`** — thin, dependency-free HTTP client for the Monit24 REST API. `client/client.go` holds the shared `Client` struct (`get`/`post`/`put`/`delete` helpers, `ResourceNotFound` error type, `OwnerID()`, and the Basic/Bearer auth selection). Each resource has its own file (`client/service.go`, `client/group.go`, `client/notification_address.go`, `client/account.go`, `client/user_data.go`, `client/user_data_setting.go`, ...) with a data struct (JSON tags matching the API) and `Create*`/`Read*`/`Update*`/`Delete*` methods.
- **`monit24/`** — the Terraform SDK provider and resources. `monit24/provider.go` defines the provider schema (`user`/`password`/`token`) and registers resources in `ResourcesMap`. Each `monit24/resource_*.go` defines the Terraform schema and `CreateContext`/`ReadContext`/`UpdateContext`/`DeleteContext` functions that translate between `*schema.ResourceData` and the corresponding `client` struct.

Key conventions to follow when touching a resource:
- IDs are API-assigned integers, stored in Terraform state as strings via `strconv.Itoa`/`strconv.Atoi`. Where there's no single numeric ID (e.g. `monit24_group_share`, `monit24_user_data_setting`), a composite `"<a>:<b>"` string ID is used instead — see `groupShareID`/`parseGroupShareID` in `monit24/resource_group_share.go` for the pattern.
- On `ReadContext`, when the API returns 404 (`client.ResourceNotFound`), call `d.SetId("")` and return `nil` (not an error) so Terraform drops it from state.
- All resources support import via `schema.ImportStatePassthroughContext`.
- Optional/nullable API fields are modeled as pointers (`*string`, `*bool`, `*[]int`, etc.) in the `client` structs, only set in `Read*` when non-nil.
- `resourceServiceCreate` delegates to `resourceServiceUpdate` after creation to populate all fields in one pass — follow this pattern for new resources with many optional fields, *unless* Update has side effects beyond a plain PUT (see `monit24_subaccount`, which deliberately does not delegate Create to Update because Update also conditionally calls the `change_password` action).
- `client.Client.OwnerID()` resolves to the parent account ID when the authenticated user is a sub-account; resources pass this as `owner_id` on create/update.
- `extended_settings` (service resource) is a free-form string map; values are coerced to int/bool/string on write (`newServiceFromResourceData`) and merged against currently-defined keys on read (`mergeMaps`) since the API can return additional settings the config doesn't declare.
- A resource whose underlying API record can't be independently deleted (it's implicitly tied to a parent's lifecycle, e.g. `monit24_user_data` — there's no `DELETE /user_data/{id}`) implements `DeleteContext` as a state-only no-op (`d.SetId(""); return nil`), not an API call.

## Resources

Eleven resources are registered in `monit24/provider.go`, tracking Monit24 API v3.51 (see `docs-internal/api-inventory.md` for the full endpoint inventory this was audited against, and `docs/superpowers/specs/2026-08-14-full-api-v3.51-coverage-design.md` for the design driving ongoing work):

- **`monit24_group`** (`resource_group.go` / `client/group.go`) — a container that other resources attach to via `group_id`. Fields: `name` (required), `periodic_daily_reports`/`periodic_weekly_reports`/`periodic_monthly_reports` (optional bool, default `true`), `archived_services_in_periodic_reports` (optional bool, default `true`), `assigned_sensor_ids` (optional set of `{category, sensor_ids}` blocks), `is_default` (computed). Does **not** expose the API's deprecated `sensor_ids` field on `group`.
- **`monit24_group_share`** (`resource_group_share.go` / `client/group_share.go`) — shares a group with another account. Composite ID `"<group_id>:<account_id>"` (no server-assigned numeric ID; `PUT /groups/{group_id}/shares/{account_id}` both creates and updates). Fields: `group_id`/`account_id` (required, `ForceNew`), five `can_*` permission bools (optional, default `false`).
- **`monit24_notification_address`** (`resource_notification_address.go` / `client/notification_address.go`) — a channel-specific address notifications are sent to. Fields: `address` (required), `notification_channel_id` (required), `group_id` (optional/computed), `description` (optional, default `""`).
- **`monit24_periodic_report_address`** (`resource_periodic_report_address.go` / `client/periodic_report_address.go`) — sibling of `notification_address` for periodic report delivery. Fields: `address` (required, email), `report_frequency` (required, `daily`/`weekly`/`monthly`), `group_id` (optional/computed).
- **`monit24_service`** (`resource_service.go` / `client/service.go`) — a monitored endpoint/check. Fields: `type_id`/`name`/`address` (required), `group_id` (optional/computed), `interval` (optional, default `600`), `description` (optional), `is_active` (optional, default `true`), `is_archived` (optional, default `false`), `sensor_ids` (optional set of ints), `step_names` (optional ordered list of strings), `notification_channel_ids`/`notification_condition_ids` (optional/computed sets of strings), `notification_mode_id`/`recovery_notification_mode_id` (optional, default `"default"`), `extended_settings` (optional/computed free-form string map). Does not expose `silent_hours`/`suspension_hours` — deprecated in favor of `weekly_suspension`. **`is_archived` write semantics are unverified against the live API** — check via `make testacc` before relying on it.
- **`monit24_subaccount`** (`resource_subaccount.go` / `client/account.go`) — creates a dependent/linked account via `POST /accounts/subaccount`. Fields: `name`/`username` (required), `package_id` (optional/computed), `is_read_only`/`disable_legacy_notifications` (optional bool, default `false`), `language_id` (optional, default `"pl"`), `time_zone_id` (optional, default `"europe_warsaw"`), `is_activated`/`is_blocked`/`parent_account_id` (computed), `subaccount_block`/`subaccount_edit`/`is_2fa_setup_required` (optional bool, default `false`, `ForceNew`), `user_data` (required, `ForceNew`, single nested block — creation-only; use `monit24_user_data` for post-creation updates), `password` (optional, sensitive, never read back, changes routed through `change_password` on Update so it doesn't force recreation), `set_password_url` (optional, `ForceNew`). **`parent_account_id`'s "inferred by the API" assumption and `DELETE /accounts/{id}` vs. `close` for removal are both unverified against the live API** — the acceptance test asserts a specific value/behavior so the next real `make testacc` run conclusively proves or disproves them.
- **`monit24_account_user`** (`resource_account_user.go`, shares `client/account.go`) — adds another login to the **same** account via plain `POST /accounts` (distinct from `monit24_subaccount`, which creates a dependent account). Same fields as `monit24_subaccount` minus the subaccount-only ones (`set_password_url`, `subaccount_block`, `subaccount_edit`, `is_2fa_setup_required`).
- **`monit24_user_data`** (`resource_user_data.go` / `client/user_data.go`) — manages the contact/billing details tied to an account (`PUT /user_data/{id}`) independently of `monit24_subaccount`/`monit24_account_user`'s creation-time `user_data` block. Fields: `account_id` (required, `ForceNew` — this *is* the resource's identity, same numeric ID as the account), `email_address` (required), `address`/`contact_person`/`phone_number`/`tax_identification_number` (optional), `ip_whitelist` (optional list of strings), `ip_whitelist_enabled` (optional bool, default `false`). `Delete` is a state-only no-op — there's no `DELETE /user_data/{id}`.
- **`monit24_user_data_setting`** (`resource_user_data_setting.go` / `client/user_data_setting.go`) — a generic per-account key/value setting (`PUT/GET/DELETE /user_data/{id}/settings/{key}`). Composite ID `"<account_id>:<key>"`. `value` is treated as a plain string — round-trips correctly for values this resource itself wrote, but can't read back settings some other client wrote as non-string JSON (number/bool/object).
- **`monit24_suspension`** (`resource_suspension.go` / `client/suspension.go`) — a one-off planned maintenance window for a service. Fields: `service_id` (required), `end_time` (required, ISO datetime), `start_time` (optional/computed), `only_notifications` (optional bool, default `false`), `description` (optional).
- **`monit24_weekly_suspension`** (`resource_weekly_suspension.go` / `client/weekly_suspension.go`) — a recurring maintenance window. `start_minute`/`end_minute` are required single-block nested objects (`day_of_week` 1-7, `hour` 0-23, `minute` 0-59) via `minuteOfWeekSchema()`. Also: `service_id` (required), `only_notifications` (optional bool, default `false`), `description` (optional).

### Known gaps / deliberately out of scope

A second "conditional alerting" subsystem exists in the API (`contacts`, `contact_groups`, `contact_addresses`, `events`, `escalations`, plus their suspension variants) — additive to, not a replacement for, `notification_address`; not yet modeled, see the design spec for the planned phased rollout. `report_templates` and `templates` (custom notification templates) are also not yet modeled. `reports` and `corrections` are intentionally never going to be resources (generate-once/point-in-time, no convergeable state). `sessions` is intentionally never going to be a resource (`POST /sessions` needs the same username+password the provider's Basic Auth already uses — no benefit to modeling it, as distinct from *using* an existing token for auth, which is supported). `admin/*` is out of scope (Monit24-staff-only API surface). `Provider().DataSourcesMap` is still empty — the design spec's Phase 1 covers dictionary data sources, Phase 6 covers per-resource companion data sources.
```

- [ ] **Step 5: Mark conflicts resolved and complete the cherry-pick**

Run: `git add monit24/provider.go CLAUDE.md && git cherry-pick --continue`
(If your editor opens for the commit message, keep the default `fb15f07` message as-is and save/close.)

- [ ] **Step 6: Verify**

Run: `gofmt -l -w . && go build ./... && go vet ./... && go test ./...`
Expected: no gofmt output (nothing needed reformatting), build/vet clean, `go test ./monit24/...` shows `ok` with all `TestAcc*` tests `SKIP` (no `TF_ACC` set).

Run: `git log --oneline -3`
Expected: the cherry-picked commit now on top of `b157cee`, both on `feature/monit24-api-v3.51-resources`.

---

### Task 2: Bearer token authentication in `client/client.go`

**Files:**
- Modify: `client/client.go`
- Create: `client/client_test.go`

**Interfaces:**
- Produces: `Client.token` field; `NewTokenClient(ctx context.Context, token string) (Client, error)`; `authorizationHeaderValue(basicAuth, token string) string` (unexported, package-level pure function).
- Consumes: nothing new from other tasks.

- [ ] **Step 1: Write the failing unit test**

Create `client/client_test.go`:

```go
package client

import "testing"

func TestAuthorizationHeaderValue(t *testing.T) {
	tests := []struct {
		name      string
		basicAuth string
		token     string
		want      string
	}{
		{
			name:      "token takes priority when both are set",
			basicAuth: "dXNlcjpwYXNz",
			token:     "abc123",
			want:      "Bearer abc123",
		},
		{
			name:      "basic auth used when no token",
			basicAuth: "dXNlcjpwYXNz",
			token:     "",
			want:      "Basic dXNlcjpwYXNz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := authorizationHeaderValue(tt.basicAuth, tt.token)
			if got != tt.want {
				t.Errorf("authorizationHeaderValue(%q, %q) = %q, want %q", tt.basicAuth, tt.token, got, tt.want)
			}
		})
	}
}
```

- [ ] **Step 2: Run it to confirm it fails to compile (function doesn't exist yet)**

Run: `go test ./client/... -run TestAuthorizationHeaderValue -v`
Expected: FAIL — `undefined: authorizationHeaderValue`

- [ ] **Step 3: Implement the auth refactor in `client/client.go`**

Replace the `Client` struct and the two constructors:

```go
type Client struct {
	client    *http.Client
	basicAuth string
	token     string
	ownerID   int
}

func NewBasicAuthClient(ctx context.Context, user string, password string) (Client, error) {
	return newAuthenticatedClient(ctx, Client{basicAuth: basicAuth(user, password), client: &http.Client{}})
}

func NewTokenClient(ctx context.Context, token string) (Client, error) {
	return newAuthenticatedClient(ctx, Client{token: token, client: &http.Client{}})
}

func newAuthenticatedClient(ctx context.Context, c Client) (Client, error) {
	account, err := c.getMyAccount(ctx)
	if err != nil {
		return Client{}, err
	}

	if account.ParentAccountID != nil {
		c.ownerID = *account.ParentAccountID
	} else {
		c.ownerID = account.ID
	}

	return c, nil
}
```

Replace the header line in `rawRequest`:

```go
	req.Header.Add("Authorization", authorizationHeaderValue(c.basicAuth, c.token))
```

Add the new pure function (near `basicAuth` at the bottom of the file):

```go
func authorizationHeaderValue(basicAuth, token string) string {
	if token != "" {
		return "Bearer " + token
	}
	return "Basic " + basicAuth
}
```

- [ ] **Step 4: Run the unit test to confirm it passes**

Run: `go test ./client/... -run TestAuthorizationHeaderValue -v`
Expected: PASS for both subtests.

- [ ] **Step 5: Full verify**

Run: `gofmt -l -w . && go build ./... && go vet ./... && go test ./...`
Expected: clean.

- [ ] **Step 6: Commit**

```bash
git add client/client.go client/client_test.go
git commit -m "Add Bearer token authentication support to client

Basic Auth (user:pass) is rejected by the API for accounts with 2FA
enabled. NewTokenClient authenticates via Authorization: Bearer <token>
instead, verified against the live API. authorizationHeaderValue picks
the right header; token takes priority when both are configured."
```

---

### Task 3: Wire `MONIT24_TOKEN` into the provider

**Files:**
- Modify: `monit24/provider.go`

**Interfaces:**
- Consumes: `client.NewTokenClient` (Task 2).
- Produces: provider schema field `token`.

- [ ] **Step 1: Add the schema field**

In `Provider()`'s `Schema` map, after the `password` field:

```go
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("MONIT24_TOKEN", nil),
			},
```

- [ ] **Step 2: Update `providerConfigure` precedence**

Replace the body of `providerConfigure` with:

```go
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := d.Get("token").(string)
	user := d.Get("user").(string)
	password := d.Get("password").(string)

	var diags diag.Diagnostics

	if token != "" {
		c, err := client.NewTokenClient(ctx, token)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return c, diags
	}

	if user != "" && password != "" {
		c, err := client.NewBasicAuthClient(ctx, user, password)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return c, diags
	}

	return client.Client{}, diag.FromErr(errors.New("no credentials provided"))
}
```

- [ ] **Step 3: Verify**

Run: `gofmt -l -w . && go build ./... && go vet ./... && go test ./...`
Expected: clean.

- [ ] **Step 4: Update `README.md`**

In `README.md`'s "## Authentication" section, after the existing `MONIT24_USER`/`MONIT24_PASSWORD` example, add:

```markdown
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
```

- [ ] **Step 5: Commit**

```bash
git add monit24/provider.go README.md
git commit -m "Add MONIT24_TOKEN provider auth mode for 2FA-enabled accounts

token takes priority over user/password when set. Documented in README
alongside the existing Basic Auth instructions."
```

---

### Task 4: `monit24_service.is_archived`

**Files:**
- Modify: `client/service.go`, `monit24/resource_service.go`

**Interfaces:**
- Consumes: `boolPtr` (existing helper in `monit24/resource_service.go`).
- Produces: `client.Service.IsArchived *bool`.

- [ ] **Step 1: Add the field to `client/service.go`**

In the `Service` struct, after `IsActive`:

```go
	IsArchived                 *bool                   `json:"is_archived,omitempty"`
```

- [ ] **Step 2: Add the schema field to `monit24/resource_service.go`**

In `resourceService()`'s `Schema`, after `"is_active"`:

```go
			"is_archived": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
```

- [ ] **Step 3: Wire it in `newServiceFromResourceData`**

After the existing `isActive := d.Get("is_active"); service.IsActive = boolPtr(isActive.(bool))` block, add:

```go
	service.IsArchived = boolPtr(d.Get("is_archived").(bool))
```

- [ ] **Step 4: Wire it in `resourceServiceRead`**

After the `if service.IsActive != nil { ... }` block, add:

```go
	if service.IsArchived != nil {
		if err := d.Set("is_archived", *service.IsArchived); err != nil {
			return diag.FromErr(err)
		}
	}
```

- [ ] **Step 5: Add an acceptance-test assertion**

In `monit24/resource_service_test.go`, add `"is_archived": "false"` to the `testServiceAttributesCreated` map and `"is_archived": "false"` to `testServiceAttributesDefaultsUpdated` (both existing `map[string]string` vars near the bottom of the file) — this doesn't change test behavior locally (acceptance tests still skip without `TF_ACC`), but documents the expected default and will catch a write-semantics regression the next time `make testacc` runs for real.

- [ ] **Step 6: Verify**

Run: `gofmt -l -w . && go build ./... && go vet ./... && go test ./...`
Expected: clean.

- [ ] **Step 7: Commit**

```bash
git add client/service.go monit24/resource_service.go monit24/resource_service_test.go
git commit -m "Add is_archived field to monit24_service

Previously excluded because the API rejected it as unsupported; it's now
a normal settable field in service_create_data per API v3.51. Write
semantics against a live account are still unverified — flagged in
CLAUDE.md."
```

---

### Task 5: Harden the `monit24_subaccount` `parent_account_id` acceptance test

**Files:**
- Modify: `monit24/resource_subaccount_test.go`

**Interfaces:**
- Consumes: `client.Client.OwnerID()` (existing).

- [ ] **Step 1: Replace the weak assertion with a precise one**

In `monit24/resource_subaccount_test.go`, find this line in the first `Check: resource.ComposeTestCheckFunc(...)` block:

```go
					resource.TestCheckResourceAttrSet("monit24_subaccount.test", "parent_account_id"),
```

Replace it with:

```go
					func(state *terraform.State) error {
						c, err := client.NewBasicAuthClient(context.Background(), os.Getenv("MONIT24_USER"), os.Getenv("MONIT24_PASSWORD"))
						if err != nil {
							return err
						}

						sub := state.RootModule().Resources["monit24_subaccount.test"]
						gotParentID := sub.Primary.Attributes["parent_account_id"]
						wantParentID := strconv.Itoa(c.OwnerID())

						if gotParentID != wantParentID {
							return fmt.Errorf("expected parent_account_id to equal the caller's own account id (%s) since it's assumed to be inferred by POST /accounts/subaccount, got %s — if this fails, parent_account_id is NOT auto-inferred; make it an Optional (not Computed-only) field in resource_subaccount.go and set it explicitly from OwnerID()", wantParentID, gotParentID)
						}

						return nil
					},
```

- [ ] **Step 2: Verify imports are already sufficient**

`context`, `os`, `strconv`, `fmt`, and `github.com/monit24/terraform-provider-monit24/client` are already imported in this file (used elsewhere in `testAccCheckSubaccountDestroyed` and the id-capture closure) — no import changes needed.

- [ ] **Step 3: Verify**

Run: `gofmt -l -w . && go build ./... && go vet ./... && go test ./...`
Expected: clean (test still `SKIP`s locally without `TF_ACC`).

- [ ] **Step 4: Commit**

```bash
git add monit24/resource_subaccount_test.go
git commit -m "Harden parent_account_id acceptance test assertion

Replaces the weak TestCheckResourceAttrSet with a precise comparison
against the caller's own OwnerID(), so the next real make testacc run
conclusively proves or disproves the 'inferred by the API' assumption
instead of just checking the field is non-empty."
```

---

### Task 6: `monit24_user_data` resource

**Files:**
- Modify: `client/account.go` (extend `UserData` struct)
- Create: `client/user_data.go`, `monit24/resource_user_data.go`, `monit24/resource_user_data_test.go`, `examples/resources/monit24_user_data/resource.tf`, `examples/resources/monit24_user_data/import.sh`
- Modify: `monit24/provider.go` (register)

**Interfaces:**
- Consumes: `client.UserData` (existing, from `client/account.go`), `strPtr`/`boolPtr` (existing helpers).
- Produces: `client.ReadUserData(ctx, id int) (UserData, error)`, `client.UpdateUserData(ctx, id int, req UserData) error`.

- [ ] **Step 1: Extend `client.UserData` with the read-only fields the full `user_data` API object has**

In `client/account.go`, replace the `UserData` struct with:

```go
type UserData struct {
	EmailAddress            string    `json:"email_address"`
	Address                 *string   `json:"address,omitempty"`
	ContactPerson           *string   `json:"contact_person,omitempty"`
	PhoneNumber             *string   `json:"phone_number,omitempty"`
	TaxIdentificationNumber *string   `json:"tax_identification_number,omitempty"`
	IPWhitelist             *[]string `json:"ip_whitelist,omitempty"`
	IPWhitelistEnabled      *bool     `json:"ip_whitelist_enabled,omitempty"`
	ID                      *int      `json:"id,omitempty"`
	CreatedAt               *string   `json:"created_at,omitempty"`
	Has2FAEnabled           *bool     `json:"has_2fa_enabled,omitempty"`
	Settings                *[]string `json:"settings,omitempty"`
}
```

(The four new fields are `omitempty` and never populated when building a write payload from Terraform config, so this doesn't change what `monit24_subaccount`/`monit24_account_user` send on create.)

- [ ] **Step 2: Create `client/user_data.go`**

```go
package client

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c Client) ReadUserData(ctx context.Context, id int) (UserData, error) {
	resp, err := c.get(ctx, fmt.Sprintf("/user_data/%v", id))
	if err != nil {
		return UserData{}, err
	}

	var userData UserData
	err = json.Unmarshal(resp, &userData)
	if err != nil {
		return UserData{}, err
	}

	return userData, nil
}

func (c Client) UpdateUserData(ctx context.Context, id int, req UserData) error {
	return c.put(ctx, fmt.Sprintf("/user_data/%v", id), req)
}
```

- [ ] **Step 3: Create `monit24/resource_user_data.go`**

```go
package monit24

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/monit24/terraform-provider-monit24/client"
)

func resourceUserData() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserDataCreate,
		ReadContext:   resourceUserDataRead,
		UpdateContext: resourceUserDataUpdate,
		DeleteContext: resourceUserDataDelete,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"email_address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"contact_person": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"phone_number": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tax_identification_number": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_whitelist": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ip_whitelist_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func flatUserDataFromResourceData(d *schema.ResourceData) client.UserData {
	userData := client.UserData{
		EmailAddress: d.Get("email_address").(string),
	}

	if v, ok := d.GetOk("address"); ok {
		userData.Address = strPtr(v.(string))
	}

	if v, ok := d.GetOk("contact_person"); ok {
		userData.ContactPerson = strPtr(v.(string))
	}

	if v, ok := d.GetOk("phone_number"); ok {
		userData.PhoneNumber = strPtr(v.(string))
	}

	if v, ok := d.GetOk("tax_identification_number"); ok {
		userData.TaxIdentificationNumber = strPtr(v.(string))
	}

	if v, ok := d.GetOk("ip_whitelist"); ok {
		list := v.([]interface{})
		ips := make([]string, len(list))

		for i := range list {
			ips[i] = list[i].(string)
		}

		userData.IPWhitelist = &ips
	}

	userData.IPWhitelistEnabled = boolPtr(d.Get("ip_whitelist_enabled").(bool))

	return userData
}

func resourceUserDataCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Client)

	accountID := d.Get("account_id").(int)
	userData := flatUserDataFromResourceData(d)

	err := c.UpdateUserData(ctx, accountID, userData)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(accountID))

	return resourceUserDataRead(ctx, d, m)
}

func resourceUserDataRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(client.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	userData, err := c.ReadUserData(ctx, id)
	if err != nil {
		if _, ok := err.(client.ResourceNotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err := d.Set("account_id", id); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("email_address", userData.EmailAddress); err != nil {
		return diag.FromErr(err)
	}

	if userData.Address != nil {
		if err := d.Set("address", *userData.Address); err != nil {
			return diag.FromErr(err)
		}
	}

	if userData.ContactPerson != nil {
		if err := d.Set("contact_person", *userData.ContactPerson); err != nil {
			return diag.FromErr(err)
		}
	}

	if userData.PhoneNumber != nil {
		if err := d.Set("phone_number", *userData.PhoneNumber); err != nil {
			return diag.FromErr(err)
		}
	}

	if userData.TaxIdentificationNumber != nil {
		if err := d.Set("tax_identification_number", *userData.TaxIdentificationNumber); err != nil {
			return diag.FromErr(err)
		}
	}

	if userData.IPWhitelist != nil {
		if err := d.Set("ip_whitelist", *userData.IPWhitelist); err != nil {
			return diag.FromErr(err)
		}
	}

	if userData.IPWhitelistEnabled != nil {
		if err := d.Set("ip_whitelist_enabled", *userData.IPWhitelistEnabled); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceUserDataUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	userData := flatUserDataFromResourceData(d)

	err = c.UpdateUserData(ctx, id, userData)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceUserDataRead(ctx, d, m)
}

func resourceUserDataDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	d.SetId("")

	return diags
}
```

- [ ] **Step 4: Register in `monit24/provider.go`**

Add to `ResourcesMap` (alphabetically, after `monit24_suspension`... actually before it, `user_data` < `weekly_suspension` but > `suspension`; insert to keep the map alphabetically sorted):

```go
			"monit24_user_data":               resourceUserData(),
```

- [ ] **Step 5: Write the acceptance test — `monit24/resource_user_data_test.go`**

```go
package monit24

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUserData(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_user_data.test", "email_address", "tf-acceptance-tests@example.com"),
					resource.TestCheckResourceAttr("monit24_user_data.test", "contact_person", "Original Contact"),
				),
			},
			{
				Config: testAccUserDataConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_user_data.test", "email_address", "tf-acceptance-tests@example.com"),
					resource.TestCheckResourceAttr("monit24_user_data.test", "contact_person", "Updated Contact"),
				),
			},
		},
	})
}

const testAccUserDataConfig = `
resource "monit24_subaccount" "test" {
  name     = "TF acceptance test user_data subaccount"
  username = "tf-acc-test-user-data"

  user_data {
    email_address = "tf-acceptance-tests@example.com"
  }
}

resource "monit24_user_data" "test" {
  account_id     = monit24_subaccount.test.id
  email_address  = "tf-acceptance-tests@example.com"
  contact_person = "Original Contact"
}
`

const testAccUserDataConfigUpdated = `
resource "monit24_subaccount" "test" {
  name     = "TF acceptance test user_data subaccount"
  username = "tf-acc-test-user-data"

  user_data {
    email_address = "tf-acceptance-tests@example.com"
  }
}

resource "monit24_user_data" "test" {
  account_id     = monit24_subaccount.test.id
  email_address  = "tf-acceptance-tests@example.com"
  contact_person = "Updated Contact"
}
`
```

- [ ] **Step 6: Add examples**

`examples/resources/monit24_user_data/resource.tf`:

```tf
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
```

`examples/resources/monit24_user_data/import.sh`:

```bash
terraform import monit24_user_data.client_a 123456
```

- [ ] **Step 7: Verify**

Run: `gofmt -l -w . && go build ./... && go vet ./... && go test ./...`
Expected: clean, new `TestAccUserData` shows `SKIP` without `TF_ACC`.

- [ ] **Step 8: Commit**

```bash
git add client/account.go client/user_data.go monit24/resource_user_data.go monit24/resource_user_data_test.go monit24/provider.go examples/resources/monit24_user_data/
git commit -m "Add monit24_user_data resource

Manages an account's contact/billing details (PUT /user_data/{id})
independently of monit24_subaccount's creation-time user_data block,
which is ForceNew because the API has no way to update it there. Delete
is a state-only no-op since there's no DELETE /user_data/{id} — the
record is tied to the account's own lifecycle."
```

---

### Task 7: `monit24_account_user` resource

**Files:**
- Modify: `client/account.go` (add `CreateAccountUser`), `monit24/provider.go` (register)
- Create: `monit24/resource_account_user.go`, `monit24/resource_account_user_test.go`, `examples/resources/monit24_account_user/resource.tf`, `examples/resources/monit24_account_user/import.sh`

**Interfaces:**
- Consumes: `client.SubaccountCreateRequest`, `client.Account`, `client.ReadAccount`/`UpdateAccount`/`ChangeAccountPassword`/`DeleteAccount` (all existing, from `client/account.go`); `accountFromResourceData`, `userDataFromResourceData` (existing package-level helpers in `monit24/resource_subaccount.go` — reused as-is, not duplicated).
- Produces: `client.CreateAccountUser(ctx, req SubaccountCreateRequest) (int, error)`.

- [ ] **Step 1: Add `CreateAccountUser` to `client/account.go`**

After `CreateSubaccount`:

```go
func (c Client) CreateAccountUser(ctx context.Context, req SubaccountCreateRequest) (int, error) {
	resp, err := c.post(ctx, "/accounts", req)
	if err != nil {
		return 0, err
	}

	var response CreateSubaccountResponse
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return 0, err
	}

	return response.ID, err
}
```

- [ ] **Step 2: Create `monit24/resource_account_user.go`**

This mirrors `resource_subaccount.go` minus the subaccount-only fields (`set_password_url`, `subaccount_block`, `subaccount_edit`, `is_2fa_setup_required`), reusing its `accountFromResourceData`/`userDataFromResourceData` helpers directly (same package):

```go
package monit24

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/monit24/terraform-provider-monit24/client"
)

func resourceAccountUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccountUserCreate,
		ReadContext:   resourceAccountUserRead,
		UpdateContext: resourceAccountUserUpdate,
		DeleteContext: resourceAccountUserDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"package_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"is_read_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"disable_legacy_notifications": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"language_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "pl",
			},
			"time_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "europe_warsaw",
			},
			"is_activated": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_blocked": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"user_data": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email_address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"address": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"contact_person": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"phone_number": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"tax_identification_number": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ip_whitelist": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"ip_whitelist_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceAccountUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Client)

	req := client.SubaccountCreateRequest{
		Account:  accountFromResourceData(d),
		UserData: userDataFromResourceData(d),
	}

	if v, ok := d.GetOk("password"); ok {
		req.Password = strPtr(v.(string))
	}

	id, err := c.CreateAccountUser(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(id))

	return resourceAccountUserRead(ctx, d, m)
}

func resourceAccountUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(client.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	account, err := c.ReadAccount(ctx, id)
	if err != nil {
		if _, ok := err.(client.ResourceNotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err := d.Set("name", account.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("username", account.Username); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("package_id", account.PackageID); err != nil {
		return diag.FromErr(err)
	}

	if account.IsReadOnly != nil {
		if err := d.Set("is_read_only", *account.IsReadOnly); err != nil {
			return diag.FromErr(err)
		}
	}

	if account.DisableLegacyNotifications != nil {
		if err := d.Set("disable_legacy_notifications", *account.DisableLegacyNotifications); err != nil {
			return diag.FromErr(err)
		}
	}

	if account.LanguageID != nil {
		if err := d.Set("language_id", *account.LanguageID); err != nil {
			return diag.FromErr(err)
		}
	}

	if account.TimeZoneID != nil {
		if err := d.Set("time_zone_id", *account.TimeZoneID); err != nil {
			return diag.FromErr(err)
		}
	}

	if account.IsActivated != nil {
		if err := d.Set("is_activated", *account.IsActivated); err != nil {
			return diag.FromErr(err)
		}
	}

	if account.IsBlocked != nil {
		if err := d.Set("is_blocked", *account.IsBlocked); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceAccountUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	account := accountFromResourceData(d)

	err = c.UpdateAccount(ctx, id, account)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("password") {
		if password := d.Get("password").(string); password != "" {
			err = c.ChangeAccountPassword(ctx, id, password)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceAccountUserRead(ctx, d, m)
}

func resourceAccountUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(client.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteAccount(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
```

- [ ] **Step 3: Register in `monit24/provider.go`**

Add to `ResourcesMap` (alphabetically, right after `monit24_group_share` and before `monit24_group`... `account_user` < `group`, so insert as the first entry):

```go
			"monit24_account_user":            resourceAccountUser(),
```

- [ ] **Step 4: Write the acceptance test — `monit24/resource_account_user_test.go`**

```go
package monit24

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAccountUser(t *testing.T) {
	username := fmt.Sprintf("tf-acc-test-user-%d", time.Now().UnixNano())

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAccountUserConfig(username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_account_user.test", "name", "TF acceptance test account user"),
					resource.TestCheckResourceAttr("monit24_account_user.test", "username", username),
				),
			},
		},
	})
}

func testAccAccountUserConfig(username string) string {
	return fmt.Sprintf(`
resource "monit24_account_user" "test" {
  name     = "TF acceptance test account user"
  username = %q

  user_data {
    email_address = "tf-acceptance-tests@example.com"
  }
}
`, username)
}
```

- [ ] **Step 5: Add examples**

`examples/resources/monit24_account_user/resource.tf`:

```tf
# A second login for the SAME account (distinct from monit24_subaccount,
# which creates a separate, dependent account).
resource "monit24_account_user" "colleague" {
  name     = "Jane's colleague"
  username = "jane-colleague"

  user_data {
    email_address = "colleague@example.com"
  }
}
```

`examples/resources/monit24_account_user/import.sh`:

```bash
terraform import monit24_account_user.colleague 123456
```

- [ ] **Step 6: Verify**

Run: `gofmt -l -w . && go build ./... && go vet ./... && go test ./...`
Expected: clean, new `TestAccAccountUser` shows `SKIP` without `TF_ACC`.

- [ ] **Step 7: Commit**

```bash
git add client/account.go monit24/resource_account_user.go monit24/resource_account_user_test.go monit24/provider.go examples/resources/monit24_account_user/
git commit -m "Add monit24_account_user resource

POST /accounts adds another login to the same account, distinct from
monit24_subaccount's POST /accounts/subaccount which creates a
dependent account. Reuses client.SubaccountCreateRequest and the
accountFromResourceData/userDataFromResourceData helpers already
defined for monit24_subaccount."
```

---

### Task 8: `monit24_user_data_setting` resource

**Files:**
- Create: `client/user_data_setting.go`, `monit24/resource_user_data_setting.go`, `monit24/resource_user_data_setting_test.go`, `examples/resources/monit24_user_data_setting/resource.tf`, `examples/resources/monit24_user_data_setting/import.sh`
- Modify: `monit24/provider.go` (register)

**Interfaces:**
- Consumes: nothing new from other tasks (standalone).
- Produces: `client.ReadUserDataSetting(ctx, id int, key string) (string, error)`, `client.PutUserDataSetting(ctx, id int, key string, value string) error`, `client.DeleteUserDataSetting(ctx, id int, key string) error`.

- [ ] **Step 1: Create `client/user_data_setting.go`**

```go
package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// The API stores settings as arbitrary JSON (string, number, bool, object,
// array). This client only ever writes plain JSON strings, so round-tripping
// a value this resource itself wrote is safe — reading a setting some other
// client wrote as a non-string JSON value will fail to unmarshal.
func (c Client) ReadUserDataSetting(ctx context.Context, id int, key string) (string, error) {
	resp, err := c.get(ctx, fmt.Sprintf("/user_data/%v/settings/%v", id, key))
	if err != nil {
		return "", err
	}

	var value string
	err = json.Unmarshal(resp, &value)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (c Client) PutUserDataSetting(ctx context.Context, id int, key string, value string) error {
	return c.put(ctx, fmt.Sprintf("/user_data/%v/settings/%v", id, key), value)
}

func (c Client) DeleteUserDataSetting(ctx context.Context, id int, key string) error {
	return c.delete(ctx, fmt.Sprintf("/user_data/%v/settings/%v", id, key))
}
```

- [ ] **Step 2: Create `monit24/resource_user_data_setting.go`**

```go
package monit24

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/monit24/terraform-provider-monit24/client"
)

func resourceUserDataSetting() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserDataSettingCreate,
		ReadContext:   resourceUserDataSettingRead,
		UpdateContext: resourceUserDataSettingUpdate,
		DeleteContext: resourceUserDataSettingDelete,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func userDataSettingID(accountID int, key string) string {
	return fmt.Sprintf("%d:%s", accountID, key)
}

func parseUserDataSettingID(id string) (accountID int, key string, err error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("invalid user data setting id %q, expected format \"<account_id>:<key>\"", id)
	}

	accountID, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", err
	}

	return accountID, parts[1], nil
}

func resourceUserDataSettingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Client)

	accountID := d.Get("account_id").(int)
	key := d.Get("key").(string)
	value := d.Get("value").(string)

	err := c.PutUserDataSetting(ctx, accountID, key, value)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(userDataSettingID(accountID, key))

	return resourceUserDataSettingRead(ctx, d, m)
}

func resourceUserDataSettingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(client.Client)

	accountID, key, err := parseUserDataSettingID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	value, err := c.ReadUserDataSetting(ctx, accountID, key)
	if err != nil {
		if _, ok := err.(client.ResourceNotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err := d.Set("account_id", accountID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("key", key); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("value", value); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceUserDataSettingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Client)

	accountID, key, err := parseUserDataSettingID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	value := d.Get("value").(string)

	err = c.PutUserDataSetting(ctx, accountID, key, value)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceUserDataSettingRead(ctx, d, m)
}

func resourceUserDataSettingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(client.Client)

	accountID, key, err := parseUserDataSettingID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteUserDataSetting(ctx, accountID, key)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
```

- [ ] **Step 3: Register in `monit24/provider.go`**

Add to `ResourcesMap` (alphabetically, right after `monit24_user_data` and before `monit24_weekly_suspension`):

```go
			"monit24_user_data_setting":       resourceUserDataSetting(),
```

- [ ] **Step 4: Write the acceptance test — `monit24/resource_user_data_setting_test.go`**

```go
package monit24

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUserDataSetting(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSettingConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_user_data_setting.test", "key", "tf_acceptance_test_setting"),
					resource.TestCheckResourceAttr("monit24_user_data_setting.test", "value", "original-value"),
				),
			},
			{
				Config: testAccUserDataSettingConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_user_data_setting.test", "value", "updated-value"),
				),
			},
		},
	})
}

const testAccUserDataSettingConfig = `
resource "monit24_subaccount" "test" {
  name     = "TF acceptance test user_data_setting subaccount"
  username = "tf-acc-test-uds"

  user_data {
    email_address = "tf-acceptance-tests@example.com"
  }
}

resource "monit24_user_data_setting" "test" {
  account_id = monit24_subaccount.test.id
  key        = "tf_acceptance_test_setting"
  value      = "original-value"
}
`

const testAccUserDataSettingConfigUpdated = `
resource "monit24_subaccount" "test" {
  name     = "TF acceptance test user_data_setting subaccount"
  username = "tf-acc-test-uds"

  user_data {
    email_address = "tf-acceptance-tests@example.com"
  }
}

resource "monit24_user_data_setting" "test" {
  account_id = monit24_subaccount.test.id
  key        = "tf_acceptance_test_setting"
  value      = "updated-value"
}
`
```

- [ ] **Step 5: Add examples**

`examples/resources/monit24_user_data_setting/resource.tf`:

```tf
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
```

`examples/resources/monit24_user_data_setting/import.sh`:

```bash
terraform import monit24_user_data_setting.dashboard_theme 123456:dashboard_theme
```

- [ ] **Step 6: Verify**

Run: `gofmt -l -w . && go build ./... && go vet ./... && go test ./...`
Expected: clean, new `TestAccUserDataSetting` shows `SKIP` without `TF_ACC`.

- [ ] **Step 7: Commit**

```bash
git add client/user_data_setting.go monit24/resource_user_data_setting.go monit24/resource_user_data_setting_test.go monit24/provider.go examples/resources/monit24_user_data_setting/
git commit -m "Add monit24_user_data_setting resource

Generic per-account key/value setting (PUT/GET/DELETE
/user_data/{id}/settings/{key}). Composite \"<account_id>:<key>\" id,
same pattern as monit24_group_share. value is treated as a plain
string — documented limitation for settings written by other clients
as non-string JSON."
```

---

### Task 9: Final docs pass and full verification

**Files:**
- Modify: `CLAUDE.md` (resource count, confirm accuracy after Tasks 4-8 landed)

**Interfaces:** none (documentation + verification only).

- [ ] **Step 1: Update the resource count and gaps section in `CLAUDE.md`**

The "## Resources" heading currently says "Eleven resources" (written in Task 1, before Tasks 6-8 added 3 more). Update it to:

```markdown
## Resources

Fourteen resources are registered in `monit24/provider.go`, tracking Monit24 API v3.51 (see `docs-internal/api-inventory.md` for the full endpoint inventory this was audited against, and `docs/superpowers/specs/2026-08-14-full-api-v3.51-coverage-design.md` for the design driving ongoing work):
```

Then insert bullets for `monit24_account_user`, `monit24_user_data`, and `monit24_user_data_setting` (content already given in Task 1 Step 4 above — they should already be present verbatim since Task 1 wrote the full target `CLAUDE.md` content in one pass; this step is just confirming nothing drifted and fixing the resource count number specifically, since Tasks 6-8 added resources after Task 1's `CLAUDE.md` rewrite).

- [ ] **Step 2: Full verification**

Run: `gofmt -l . && echo "gofmt clean" && go build ./... && echo "build OK" && go vet ./... && echo "vet OK" && go test ./... -v 2>&1 | tail -40`
Expected: no gofmt output, build/vet OK, all `TestAcc*` tests (including the 3 new ones from Tasks 6-8) show `SKIP`, overall `ok`.

Note: `make docs` is expected to still fail on this machine (`tfenv` has no resolvable default Terraform version and `tfenv use` itself fails without GNU grep) — this is a pre-existing local environment gap, not something this plan's tasks fix. Don't attempt `make docs` as part of verification; don't silently claim docs are regenerated. `docs/resources/*.md` for the new resources will need generating on a machine/CI with a working `tfenv`/`terraform` setup before merge.

- [ ] **Step 3: Review the full diff one more time**

Run: `git log --oneline main..HEAD`
Expected: one commit per task (cherry-pick + 8 feature commits), all on `feature/monit24-api-v3.51-resources`, none pushed.

- [ ] **Step 4: Commit the final `CLAUDE.md` fix (if Step 1 changed anything)**

```bash
git add CLAUDE.md
git commit -m "Update CLAUDE.md resource count after account-family additions"
```

(If Step 1 found nothing to change, skip this commit.)
