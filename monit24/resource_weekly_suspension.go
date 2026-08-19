package monit24

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/monit24/terraform-provider-monit24/client"
)

func minuteOfWeekSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"day_of_week": {
					Type:     schema.TypeInt,
					Required: true,
				},
				"hour": {
					Type:     schema.TypeInt,
					Required: true,
				},
				"minute": {
					Type:     schema.TypeInt,
					Required: true,
				},
			},
		},
	}
}

func resourceWeeklySuspension() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWeeklySuspensionCreate,
		ReadContext:   resourceWeeklySuspensionRead,
		UpdateContext: resourceWeeklySuspensionUpdate,
		DeleteContext: resourceWeeklySuspensionDelete,
		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"start_minute": minuteOfWeekSchema(),
			"end_minute":   minuteOfWeekSchema(),
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

func minuteOfWeekFromResourceData(v interface{}) client.MinuteOfWeek {
	m := v.([]interface{})[0].(map[string]interface{})

	return client.MinuteOfWeek{
		DayOfWeek: m["day_of_week"].(int),
		Hour:      m["hour"].(int),
		Minute:    m["minute"].(int),
	}
}

func minuteOfWeekToResourceData(m client.MinuteOfWeek) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"day_of_week": m.DayOfWeek,
			"hour":        m.Hour,
			"minute":      m.Minute,
		},
	}
}

func weeklySuspensionFromResourceData(d *schema.ResourceData) client.WeeklySuspension {
	suspension := client.WeeklySuspension{
		ServiceID:         d.Get("service_id").(int),
		StartMinute:       minuteOfWeekFromResourceData(d.Get("start_minute")),
		EndMinute:         minuteOfWeekFromResourceData(d.Get("end_minute")),
		OnlyNotifications: boolPtr(d.Get("only_notifications").(bool)),
	}

	// HasChange, not GetOk: GetOk can't distinguish "never configured" from
	// "explicitly cleared" (both read as ""), so clearing a previously-set
	// description would never reach the API.
	if d.HasChange("description") {
		suspension.Description = strPtr(d.Get("description").(string))
	}

	return suspension
}

func resourceWeeklySuspensionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Client)

	suspension := weeklySuspensionFromResourceData(d)

	id, err := c.CreateWeeklySuspension(ctx, suspension)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(id))

	return resourceWeeklySuspensionRead(ctx, d, m)
}

func resourceWeeklySuspensionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(client.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	suspension, err := c.ReadWeeklySuspension(ctx, id)
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

	if err := d.Set("start_minute", minuteOfWeekToResourceData(suspension.StartMinute)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("end_minute", minuteOfWeekToResourceData(suspension.EndMinute)); err != nil {
		return diag.FromErr(err)
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

func resourceWeeklySuspensionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Client)

	suspension := weeklySuspensionFromResourceData(d)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.UpdateWeeklySuspension(ctx, id, suspension)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceWeeklySuspensionRead(ctx, d, m)
}

func resourceWeeklySuspensionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(client.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteWeeklySuspension(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
