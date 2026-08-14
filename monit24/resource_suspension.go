package monit24

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/monit24/terraform-provider-monit24/client"
)

func resourceSuspension() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSuspensionCreate,
		ReadContext:   resourceSuspensionRead,
		UpdateContext: resourceSuspensionUpdate,
		DeleteContext: resourceSuspensionDelete,
		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"only_notifications": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func suspensionFromResourceData(d *schema.ResourceData) client.Suspension {
	suspension := client.Suspension{
		ServiceID:         d.Get("service_id").(int),
		EndTime:           d.Get("end_time").(string),
		OnlyNotifications: boolPtr(d.Get("only_notifications").(bool)),
	}

	if v, ok := d.GetOk("start_time"); ok {
		suspension.StartTime = strPtr(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		suspension.Description = strPtr(v.(string))
	}

	return suspension
}

func resourceSuspensionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Client)

	suspension := suspensionFromResourceData(d)

	id, err := c.CreateSuspension(ctx, suspension)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(id))

	return resourceSuspensionRead(ctx, d, m)
}

func resourceSuspensionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(client.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	suspension, err := c.ReadSuspension(ctx, id)
	if err != nil {
		if _, ok := err.(client.ResourceNotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err := d.Set("service_id", suspension.ServiceID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("end_time", suspension.EndTime); err != nil {
		return diag.FromErr(err)
	}

	if suspension.StartTime != nil {
		if err := d.Set("start_time", *suspension.StartTime); err != nil {
			return diag.FromErr(err)
		}
	}

	if suspension.OnlyNotifications != nil {
		if err := d.Set("only_notifications", *suspension.OnlyNotifications); err != nil {
			return diag.FromErr(err)
		}
	}

	if suspension.Description != nil {
		if err := d.Set("description", *suspension.Description); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceSuspensionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Client)

	suspension := suspensionFromResourceData(d)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.UpdateSuspension(ctx, id, suspension)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSuspensionRead(ctx, d, m)
}

func resourceSuspensionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(client.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteSuspension(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
