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
