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
