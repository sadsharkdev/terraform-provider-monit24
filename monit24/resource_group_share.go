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

func resourceGroupShare() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupShareCreate,
		ReadContext:   resourceGroupShareRead,
		UpdateContext: resourceGroupShareUpdate,
		DeleteContext: resourceGroupShareDelete,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"account_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"can_archive_services": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"can_create_services": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"can_delete_services": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"can_force_analyses": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"can_modify_group": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"can_modify_services": {
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

func groupShareID(groupID, accountID int) string {
	return fmt.Sprintf("%d:%d", groupID, accountID)
}

func parseGroupShareID(id string) (groupID int, accountID int, err error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid group share id %q, expected format \"<group_id>:<account_id>\"", id)
	}

	groupID, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	accountID, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}

	return groupID, accountID, nil
}

func groupShareFromResourceData(d *schema.ResourceData) client.GroupShare {
	return client.GroupShare{
		GroupID:            d.Get("group_id").(int),
		AccountID:          d.Get("account_id").(int),
		CanArchiveServices: boolPtr(d.Get("can_archive_services").(bool)),
		CanCreateServices:  boolPtr(d.Get("can_create_services").(bool)),
		CanDeleteServices:  boolPtr(d.Get("can_delete_services").(bool)),
		CanForceAnalyses:   boolPtr(d.Get("can_force_analyses").(bool)),
		CanModifyGroup:     boolPtr(d.Get("can_modify_group").(bool)),
		CanModifyServices:  boolPtr(d.Get("can_modify_services").(bool)),
	}
}

func resourceGroupShareCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Client)

	groupID := d.Get("group_id").(int)
	accountID := d.Get("account_id").(int)
	share := groupShareFromResourceData(d)

	err := c.PutGroupShare(ctx, groupID, accountID, share)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(groupShareID(groupID, accountID))

	return resourceGroupShareRead(ctx, d, m)
}

func resourceGroupShareRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(client.Client)

	groupID, accountID, err := parseGroupShareID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	share, err := c.ReadGroupShare(ctx, groupID, accountID)
	if err != nil {
		if _, ok := err.(client.ResourceNotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err := d.Set("group_id", share.GroupID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("account_id", share.AccountID); err != nil {
		return diag.FromErr(err)
	}

	if share.CanArchiveServices != nil {
		if err := d.Set("can_archive_services", *share.CanArchiveServices); err != nil {
			return diag.FromErr(err)
		}
	}

	if share.CanCreateServices != nil {
		if err := d.Set("can_create_services", *share.CanCreateServices); err != nil {
			return diag.FromErr(err)
		}
	}

	if share.CanDeleteServices != nil {
		if err := d.Set("can_delete_services", *share.CanDeleteServices); err != nil {
			return diag.FromErr(err)
		}
	}

	if share.CanForceAnalyses != nil {
		if err := d.Set("can_force_analyses", *share.CanForceAnalyses); err != nil {
			return diag.FromErr(err)
		}
	}

	if share.CanModifyGroup != nil {
		if err := d.Set("can_modify_group", *share.CanModifyGroup); err != nil {
			return diag.FromErr(err)
		}
	}

	if share.CanModifyServices != nil {
		if err := d.Set("can_modify_services", *share.CanModifyServices); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceGroupShareUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Client)

	groupID, accountID, err := parseGroupShareID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	share := groupShareFromResourceData(d)

	err = c.PutGroupShare(ctx, groupID, accountID, share)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGroupShareRead(ctx, d, m)
}

func resourceGroupShareDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(client.Client)

	groupID, accountID, err := parseGroupShareID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteGroupShare(ctx, groupID, accountID)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
