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
		DeleteContext: resourceAccountUserDelete,
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
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
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

func resourceAccountUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
