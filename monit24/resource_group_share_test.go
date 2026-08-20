package monit24

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/monit24/terraform-provider-monit24/client"
)

func TestAccGroupShare(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			if os.Getenv("MONIT24_SHARE_ACCOUNT_ID") == "" {
				t.Skip("MONIT24_SHARE_ACCOUNT_ID must be set to an account ID the test account can share groups with")
			}
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupShareConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_group_share.test", "can_modify_group", "false"),
					resource.TestCheckResourceAttr("monit24_group_share.test", "can_create_services", "false"),
				),
			},
			{
				Config: testAccGroupShareConfigUpdated(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_group_share.test", "can_modify_group", "true"),
					resource.TestCheckResourceAttr("monit24_group_share.test", "can_create_services", "true"),
				),
			},
		},
	})
}

func testAccGroupShareConfig() string {
	accountID, _ := strconv.Atoi(os.Getenv("MONIT24_SHARE_ACCOUNT_ID"))

	return fmt.Sprintf(`
resource "monit24_group_share" "test" {
  group_id   = monit24_group.test.id
  account_id = %d
}

resource "monit24_group" "test" {
  name = "test group"
}
`, accountID)
}

func testAccGroupShareConfigUpdated() string {
	accountID, _ := strconv.Atoi(os.Getenv("MONIT24_SHARE_ACCOUNT_ID"))

	return fmt.Sprintf(`
resource "monit24_group_share" "test" {
  group_id              = monit24_group.test.id
  account_id            = %d
  can_modify_group      = true
  can_create_services   = true
}

resource "monit24_group" "test" {
  name = "test group"
}
`, accountID)
}

func TestGroupShareIDRoundTrip(t *testing.T) {
	id := groupShareID(123, 456)
	if id != "123:456" {
		t.Fatalf("expected %q, got %q", "123:456", id)
	}

	groupID, accountID, err := parseGroupShareID(id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if groupID != 123 || accountID != 456 {
		t.Fatalf("expected (123, 456), got (%d, %d)", groupID, accountID)
	}
}

func TestParseGroupShareIDInvalid(t *testing.T) {
	tests := []string{"", "123", "123:abc", "abc:456"}

	for _, id := range tests {
		if _, _, err := parseGroupShareID(id); err == nil {
			t.Errorf("parseGroupShareID(%q): expected an error, got nil", id)
		}
	}
}

func TestGroupShareFromResourceData(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceGroupShare().Schema, map[string]interface{}{
		"group_id":         1,
		"account_id":       2,
		"can_modify_group": true,
	})

	share := groupShareFromResourceData(d)

	if share.GroupID != 1 || share.AccountID != 2 {
		t.Errorf("expected group_id=1 account_id=2, got %+v", share)
	}
	if share.CanModifyGroup == nil || !*share.CanModifyGroup {
		t.Errorf("expected can_modify_group=true, got %v", share.CanModifyGroup)
	}
	if share.CanCreateServices == nil || *share.CanCreateServices {
		t.Errorf("expected can_create_services default false, got %v", share.CanCreateServices)
	}
}

func TestGroupShareCreateUpdateDeleteLifecycle(t *testing.T) {
	var stored client.GroupShare

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPut && r.URL.Path == "/groups/1/shares/2":
			json.NewDecoder(r.Body).Decode(&stored)
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodGet && r.URL.Path == "/groups/1/shares/2":
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(stored)
		case r.Method == http.MethodDelete && r.URL.Path == "/groups/1/shares/2":
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	c := client.NewTestClient(server.URL)

	d := schema.TestResourceDataRaw(t, resourceGroupShare().Schema, map[string]interface{}{
		"group_id":   1,
		"account_id": 2,
	})

	if diags := resourceGroupShareCreate(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error creating group_share: %v", diags)
	}
	if d.Id() != "1:2" {
		t.Fatalf("expected id \"1:2\", got %q", d.Id())
	}
	if stored.GroupID != 1 || stored.AccountID != 2 {
		t.Errorf("expected create request to carry group_id/account_id, server stored %+v", stored)
	}

	if diags := resourceGroupShareRead(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error reading group_share: %v", diags)
	}
	if d.Get("group_id").(int) != 1 {
		t.Errorf("expected group_id to round-trip through read, got %v", d.Get("group_id"))
	}

	if err := d.Set("can_modify_group", true); err != nil {
		t.Fatalf("unexpected error setting can_modify_group: %v", err)
	}
	if diags := resourceGroupShareUpdate(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error updating group_share: %v", diags)
	}
	if stored.CanModifyGroup == nil || !*stored.CanModifyGroup {
		t.Errorf("expected update request to carry can_modify_group=true, server stored %+v", stored.CanModifyGroup)
	}

	if diags := resourceGroupShareDelete(context.Background(), d, c); diags.HasError() {
		t.Fatalf("unexpected error deleting group_share: %v", diags)
	}
}
