package monit24

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/monit24/terraform-provider-monit24/client"
)

func resourceSubaccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSubaccountCreate,
		ReadContext:   resourceSubaccountRead,
		UpdateContext: resourceSubaccountUpdate,
		DeleteContext: resourceSubaccountDelete,
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
			"parent_account_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			// Not readable back from the API — Terraform trusts the last-known
			// config/state value across refreshes. Omit to create the subaccount
			// without a usable password (pair with set_password_url, or reset it
			// out-of-band later). Changing this value routes through the
			// dedicated change_password endpoint on Update, not a full PUT, so it
			// never forces recreation of the resource.
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"set_password_url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"subaccount_block": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"subaccount_edit": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"is_2fa_setup_required": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			// Not readable back from the API, and not updatable via this
			// resource (the API has a separate /user_data/{id} endpoint for
			// that) — set once at creation.
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

func accountFromResourceData(d *schema.ResourceData) client.Account {
	account := client.Account{
		Name:                       d.Get("name").(string),
		Username:                   d.Get("username").(string),
		IsReadOnly:                 boolPtr(d.Get("is_read_only").(bool)),
		DisableLegacyNotifications: boolPtr(d.Get("disable_legacy_notifications").(bool)),
		LanguageID:                 strPtr(d.Get("language_id").(string)),
		TimeZoneID:                 strPtr(d.Get("time_zone_id").(string)),
	}

	if v, ok := d.GetOk("package_id"); ok {
		account.PackageID = intPtr(v.(int))
	}

	return account
}

func userDataFromResourceData(d *schema.ResourceData) client.UserData {
	m := d.Get("user_data").([]interface{})[0].(map[string]interface{})

	userData := client.UserData{
		EmailAddress: m["email_address"].(string),
	}

	if v, ok := m["address"].(string); ok && v != "" {
		userData.Address = strPtr(v)
	}

	if v, ok := m["contact_person"].(string); ok && v != "" {
		userData.ContactPerson = strPtr(v)
	}

	if v, ok := m["phone_number"].(string); ok && v != "" {
		userData.PhoneNumber = strPtr(v)
	}

	if v, ok := m["tax_identification_number"].(string); ok && v != "" {
		userData.TaxIdentificationNumber = strPtr(v)
	}

	if v, ok := m["ip_whitelist"].([]interface{}); ok && len(v) > 0 {
		ips := make([]string, len(v))
		for i := range v {
			ips[i] = v[i].(string)
		}
		userData.IPWhitelist = &ips
	}

	if v, ok := m["ip_whitelist_enabled"].(bool); ok {
		userData.IPWhitelistEnabled = boolPtr(v)
	}

	return userData
}

func resourceSubaccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Client)

	req := client.SubaccountCreateRequest{
		Account:            accountFromResourceData(d),
		UserData:           userDataFromResourceData(d),
		SubaccountBlock:    boolPtr(d.Get("subaccount_block").(bool)),
		SubaccountEdit:     boolPtr(d.Get("subaccount_edit").(bool)),
		Is2FASetupRequired: boolPtr(d.Get("is_2fa_setup_required").(bool)),
	}

	if v, ok := d.GetOk("password"); ok {
		req.Password = strPtr(v.(string))
	}

	if v, ok := d.GetOk("set_password_url"); ok {
		req.SetPasswordURL = strPtr(v.(string))
	}

	id, err := c.CreateSubaccount(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(id))

	return resourceSubaccountRead(ctx, d, m)
}

func resourceSubaccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	if account.PackageID != nil {
		if err := d.Set("package_id", *account.PackageID); err != nil {
			return diag.FromErr(err)
		}
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

	if account.ParentAccountID != nil {
		if err := d.Set("parent_account_id", *account.ParentAccountID); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceSubaccountUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := updateAccountAndPassword(ctx, c, id, d); err != nil {
		return diag.FromErr(err)
	}

	return resourceSubaccountRead(ctx, d, m)
}

// updateAccountAndPassword is shared by monit24_subaccount and
// monit24_account_user, which both PUT the same client.Account shape and
// route password changes through the dedicated change_password action
// instead of the account PUT. Skipping UpdateAccount when only password
// changed avoids sending a redundant PUT on password-only rotations.
func updateAccountAndPassword(ctx context.Context, c client.Client, id int, d *schema.ResourceData) error {
	if d.HasChangesExcept("password") {
		account := accountFromResourceData(d)

		if err := c.UpdateAccount(ctx, id, account); err != nil {
			return err
		}
	}

	if d.HasChange("password") {
		if password := d.Get("password").(string); password != "" {
			if err := c.ChangeAccountPassword(ctx, id, password); err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceSubaccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
