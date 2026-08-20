package monit24

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/monit24/terraform-provider-monit24/client"
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

func TestUserDataSettingIDRoundTrip(t *testing.T) {
	id := userDataSettingID(123456, "dashboard_theme")
	if id != "123456:dashboard_theme" {
		t.Fatalf("expected %q, got %q", "123456:dashboard_theme", id)
	}

	accountID, key, err := parseUserDataSettingID(id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if accountID != 123456 || key != "dashboard_theme" {
		t.Fatalf("expected (123456, \"dashboard_theme\"), got (%d, %q)", accountID, key)
	}
}

func TestUserDataSettingIDKeyContainingColon(t *testing.T) {
	// SplitN(id, ":", 2) must only split on the first colon, since the key
	// itself may legitimately contain one.
	accountID, key, err := parseUserDataSettingID("123456:namespace:setting")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if accountID != 123456 || key != "namespace:setting" {
		t.Fatalf("expected (123456, \"namespace:setting\"), got (%d, %q)", accountID, key)
	}
}

func TestParseUserDataSettingIDInvalid(t *testing.T) {
	tests := []string{"", "123456", "abc:dashboard_theme"}

	for _, id := range tests {
		if _, _, err := parseUserDataSettingID(id); err == nil {
			t.Errorf("parseUserDataSettingID(%q): expected an error, got nil", id)
		}
	}
}

func TestUserDataSettingCreateUpdateDeleteLifecycle(t *testing.T) {
	var stored string
	var getCount int
	var deleted bool

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPut && r.URL.Path == "/user_data/1/settings/theme":
			body, _ := io.ReadAll(r.Body)
			var v string
			json.Unmarshal(body, &v)
			stored = v
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodGet && r.URL.Path == "/user_data/1/settings/theme":
			getCount++
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(stored)
		case r.Method == http.MethodDelete && r.URL.Path == "/user_data/1/settings/theme":
			deleted = true
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	c := client.NewTestClient(server.URL)

	d := schema.TestResourceDataRaw(t, resourceUserDataSetting().Schema, map[string]interface{}{
		"account_id": 1,
		"key":        "theme",
		"value":      "dark",
	})

	if diags := resourceUserDataSettingCreate(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error creating user_data_setting: %v", diags)
	}
	if d.Id() != "1:theme" {
		t.Fatalf("expected id \"1:theme\", got %q", d.Id())
	}
	if stored != "dark" {
		t.Errorf("expected create request to carry value=dark, server stored %q", stored)
	}
	// No Computed fields in this schema, so Create must not re-fetch via Read.
	if getCount != 0 {
		t.Errorf("expected Create to not call Read, but GET was called %d time(s)", getCount)
	}

	if diags := resourceUserDataSettingRead(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error reading user_data_setting: %v", diags)
	}
	if d.Get("value").(string) != "dark" {
		t.Errorf("expected value to round-trip through read, got %q", d.Get("value"))
	}

	if err := d.Set("value", "light"); err != nil {
		t.Fatalf("unexpected error setting value: %v", err)
	}
	if diags := resourceUserDataSettingUpdate(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error updating user_data_setting: %v", diags)
	}
	if stored != "light" {
		t.Errorf("expected update request to carry the new value, server stored %q", stored)
	}

	if diags := resourceUserDataSettingDelete(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error deleting user_data_setting: %v", diags)
	}
	if !deleted {
		t.Error("expected Delete to call DeleteUserDataSetting")
	}
}
