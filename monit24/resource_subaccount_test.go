package monit24

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/monit24/terraform-provider-monit24/client"
)

func TestAccSubaccount(t *testing.T) {
	var subaccountID int

	username := fmt.Sprintf("tf-acc-test-%d", time.Now().UnixNano())

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckSubaccountDestroyed(&subaccountID),
		Steps: []resource.TestStep{
			{
				Config: testAccSubaccountConfig(username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_subaccount.test", "name", "TF acceptance test subaccount"),
					resource.TestCheckResourceAttr("monit24_subaccount.test", "username", username),
					resource.TestCheckResourceAttr("monit24_subaccount.test", "user_data.0.email_address", "tf-acceptance-tests@example.com"),
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
					func(state *terraform.State) error {
						sub := state.RootModule().Resources["monit24_subaccount.test"]
						id, err := strconv.Atoi(sub.Primary.ID)
						if err != nil {
							return err
						}

						subaccountID = id

						return nil
					},
				),
			},
			{
				// Setting a password on an existing subaccount must go through
				// change_password, not a full PUT/recreate — the resource id
				// must stay the same as the previous step.
				Config: testAccSubaccountConfigWithPassword(username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("monit24_subaccount.test", "name", "TF acceptance test subaccount"),
					func(state *terraform.State) error {
						sub := state.RootModule().Resources["monit24_subaccount.test"]
						id, err := strconv.Atoi(sub.Primary.ID)
						if err != nil {
							return err
						}

						if id != subaccountID {
							return fmt.Errorf("expected subaccount id to stay %d after changing password, got %d (resource was recreated)", subaccountID, id)
						}

						return nil
					},
				),
			},
		},
	})
}

func testAccCheckSubaccountDestroyed(subaccountID *int) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		c, err := client.NewBasicAuthClient(context.Background(), os.Getenv("MONIT24_USER"), os.Getenv("MONIT24_PASSWORD"))
		if err != nil {
			return err
		}

		_, err = c.ReadAccount(context.Background(), *subaccountID)
		if _, ok := err.(client.ResourceNotFound); ok {
			return nil
		}
		if err != nil {
			return err
		}

		return fmt.Errorf("subaccount %d still exists after destroy", *subaccountID)
	}
}

func testAccSubaccountConfig(username string) string {
	return fmt.Sprintf(`
resource "monit24_subaccount" "test" {
  name     = "TF acceptance test subaccount"
  username = %q

  user_data {
    email_address = "tf-acceptance-tests@example.com"
  }
}
`, username)
}

func testAccSubaccountConfigWithPassword(username string) string {
	return fmt.Sprintf(`
resource "monit24_subaccount" "test" {
  name     = "TF acceptance test subaccount"
  username = %q
  password = "correct-horse-battery-staple"

  user_data {
    email_address = "tf-acceptance-tests@example.com"
  }
}
`, username)
}
