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
		DeleteContext: resourceDeleteAccount,
		Schema:        accountSchemaFields(),
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

	if err := setAccountFields(d, account); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceAccountUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := updateAccountAndPassword(ctx, c, id, d); err != nil {
		return diag.FromErr(err)
	}

	return resourceAccountUserRead(ctx, d, m)
}
