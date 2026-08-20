package monit24

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/monit24/terraform-provider-monit24/client"
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

func TestFlatUserDataFromResourceData(t *testing.T) {
	t.Run("required field only", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceUserData().Schema, map[string]interface{}{
			"account_id":    1,
			"email_address": "test@example.com",
		})

		userData := flatUserDataFromResourceData(d)

		if userData.EmailAddress != "test@example.com" {
			t.Errorf("expected email_address to round-trip, got %q", userData.EmailAddress)
		}
		if userData.Address != nil {
			t.Errorf("expected address to stay nil when not configured, got %v", *userData.Address)
		}
		if userData.IPWhitelist != nil {
			t.Errorf("expected ip_whitelist to stay nil when not configured, got %v", *userData.IPWhitelist)
		}
		if userData.IPWhitelistEnabled == nil || *userData.IPWhitelistEnabled {
			t.Errorf("expected ip_whitelist_enabled default false, got %v", userData.IPWhitelistEnabled)
		}
	})

	t.Run("optional fields and ip_whitelist set", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceUserData().Schema, map[string]interface{}{
			"account_id":           1,
			"email_address":        "test@example.com",
			"address":              "1 Example St",
			"ip_whitelist":         []interface{}{"1.2.3.4"},
			"ip_whitelist_enabled": true,
		})

		userData := flatUserDataFromResourceData(d)

		if userData.Address == nil || *userData.Address != "1 Example St" {
			t.Errorf("expected address to round-trip, got %v", userData.Address)
		}
		if userData.IPWhitelist == nil || len(*userData.IPWhitelist) != 1 || (*userData.IPWhitelist)[0] != "1.2.3.4" {
			t.Errorf("expected ip_whitelist to round-trip, got %v", userData.IPWhitelist)
		}
		if userData.IPWhitelistEnabled == nil || !*userData.IPWhitelistEnabled {
			t.Errorf("expected ip_whitelist_enabled=true, got %v", userData.IPWhitelistEnabled)
		}
	})
}

func TestUserDataCreateUpdateDeleteLifecycle(t *testing.T) {
	var stored client.UserData
	var getCount int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPut && r.URL.Path == "/user_data/1":
			if err := json.NewDecoder(r.Body).Decode(&stored); err != nil {
				t.Fatalf("failed to decode PUT body: %v", err)
			}
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodGet && r.URL.Path == "/user_data/1":
			getCount++
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(stored)
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	c := client.NewTestClient(server.URL)

	d := schema.TestResourceDataRaw(t, resourceUserData().Schema, map[string]interface{}{
		"account_id":    1,
		"email_address": "test@example.com",
	})

	if diags := resourceUserDataCreate(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error creating user_data: %v", diags)
	}
	if d.Id() != "1" {
		t.Fatalf("expected id \"1\", got %q", d.Id())
	}
	if stored.EmailAddress != "test@example.com" {
		t.Errorf("expected create request to carry email_address, server stored %+v", stored)
	}
	// There are no Computed fields in this schema, so Create must not
	// re-fetch via Read — a GET here would mean the no-redundant-Read fix
	// from the code review regressed.
	if getCount != 0 {
		t.Errorf("expected Create to not call Read (no Computed fields), but GET was called %d time(s)", getCount)
	}

	if diags := resourceUserDataRead(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error reading user_data: %v", diags)
	}
	if d.Get("email_address").(string) != "test@example.com" {
		t.Errorf("expected email_address to round-trip through read, got %q", d.Get("email_address"))
	}
	if getCount != 1 {
		t.Errorf("expected exactly 1 GET from the explicit Read call, got %d", getCount)
	}

	// d.Set doesn't affect HasChange (it writes the "set" overlay, not the
	// diff), so a fresh ResourceData built from raw config is used here to
	// get a genuine "configured" diff for contact_person, with the id
	// carried over manually — matching flatUserDataFromResourceData's
	// HasChange gate.
	d = schema.TestResourceDataRaw(t, resourceUserData().Schema, map[string]interface{}{
		"account_id":     1,
		"email_address":  "test@example.com",
		"contact_person": "Jane Doe",
	})
	d.SetId("1")

	if diags := resourceUserDataUpdate(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error updating user_data: %v", diags)
	}
	if stored.ContactPerson == nil || *stored.ContactPerson != "Jane Doe" {
		t.Errorf("expected update request to carry the new contact_person, server stored %+v", stored.ContactPerson)
	}
	if getCount != 1 {
		t.Errorf("expected Update to not call Read either, but GET count changed to %d", getCount)
	}

	if diags := resourceUserDataDelete(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error deleting user_data: %v", diags)
	}
	if d.Id() != "" {
		t.Errorf("expected Delete to be a state-only no-op clearing the id, got %q", d.Id())
	}
}
