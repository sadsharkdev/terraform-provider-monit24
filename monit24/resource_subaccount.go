package monit24

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/monit24/terraform-provider-monit24/client"
)

// accountSchemaFields returns the schema fields shared by monit24_subaccount
// and monit24_account_user. Each call builds a fresh map of fresh *schema.Schema
// values, so the two resources never alias each other's schema objects.
func accountSchemaFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
		// Not readable back from the API — Terraform trusts the last-known
		// config/state value across refreshes. Omit to create the account
		// without a usable password (pair with set_password_url, or reset it
		// out-of-band later). Changing this value routes through the
		// dedicated change_password endpoint on Update, not a full PUT, so it
		// never forces recreation of the resource.
		"password": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
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
				Schema: userDataSchemaFields(),
			},
		},
	}
}

func userDataSchemaFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
	}
}

func resourceSubaccount() *schema.Resource {
	fields := accountSchemaFields()
	fields["parent_account_id"] = &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	}
	fields["set_password_url"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
	}
	fields["subaccount_block"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
		ForceNew: true,
	}
	fields["subaccount_edit"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
		ForceNew: true,
	}
	fields["is_2fa_setup_required"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
		ForceNew: true,
	}

	return &schema.Resource{
		CreateContext: resourceSubaccountCreate,
		ReadContext:   resourceSubaccountRead,
		UpdateContext: resourceSubaccountUpdate,
		DeleteContext: resourceDeleteAccount,
		Schema:        fields,
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

// setAccountFields writes the fields shared by monit24_subaccount and
// monit24_account_user into d. Callers handle their own resource-specific
// fields (e.g. subaccount's parent_account_id) after calling this.
func setAccountFields(d *schema.ResourceData, account client.Account) error {
	if err := d.Set("name", account.Name); err != nil {
		return err
	}

	if err := d.Set("username", account.Username); err != nil {
		return err
	}

	if account.PackageID != nil {
		if err := d.Set("package_id", *account.PackageID); err != nil {
			return err
		}
	}

	if account.IsReadOnly != nil {
		if err := d.Set("is_read_only", *account.IsReadOnly); err != nil {
			return err
		}
	}

	if account.DisableLegacyNotifications != nil {
		if err := d.Set("disable_legacy_notifications", *account.DisableLegacyNotifications); err != nil {
			return err
		}
	}

	if account.LanguageID != nil {
		if err := d.Set("language_id", *account.LanguageID); err != nil {
			return err
		}
	}

	if account.TimeZoneID != nil {
		if err := d.Set("time_zone_id", *account.TimeZoneID); err != nil {
			return err
		}
	}

	if account.IsActivated != nil {
		if err := d.Set("is_activated", *account.IsActivated); err != nil {
			return err
		}
	}

	if account.IsBlocked != nil {
		if err := d.Set("is_blocked", *account.IsBlocked); err != nil {
			return err
		}
	}

	return nil
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

	if err := setAccountFields(d, account); err != nil {
		return diag.FromErr(err)
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
		password := d.Get("password").(string)
		if password == "" {
			// Silently no-op'ing here would let Terraform believe the
			// password was cleared (it's written into state regardless,
			// per the password field's own doc comment) while the real
			// account password on the server stays whatever it was —
			// permanent, undetectable drift. There's no API operation to
			// reset a password to empty, so reject the transition instead.
			return fmt.Errorf("password cannot be cleared once set (there is no API operation to reset it to empty) — set it back to its previous value, or manage password rotation via set_password_url instead")
		}

		if err := c.ChangeAccountPassword(ctx, id, password); err != nil {
			return err
		}
	}

	return nil
}

// resourceDeleteAccount is shared by monit24_subaccount and
// monit24_account_user: both just parse the id and call DeleteAccount, with
// no resource-specific behavior to diverge on.
func resourceDeleteAccount(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
