package monit24

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
