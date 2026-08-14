package monit24

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/monit24/terraform-provider-monit24/client"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"periodic_daily_reports": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"periodic_weekly_reports": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"periodic_monthly_reports": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"archived_services_in_periodic_reports": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"assigned_sensor_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:     schema.TypeString,
							Required: true,
						},
						"sensor_ids": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
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

func groupFromResourceData(d *schema.ResourceData, c client.Client) client.Group {
	group := client.Group{
		Name:                              d.Get("name").(string),
		OwnerID:                           c.OwnerID(),
		PeriodicDailyReports:              boolPtr(d.Get("periodic_daily_reports").(bool)),
		PeriodicWeeklyReports:             boolPtr(d.Get("periodic_weekly_reports").(bool)),
		PeriodicMonthlyReports:            boolPtr(d.Get("periodic_monthly_reports").(bool)),
		ArchivedServicesInPeriodicReports: boolPtr(d.Get("archived_services_in_periodic_reports").(bool)),
	}

	set := d.Get("assigned_sensor_ids").(*schema.Set).List()
	assigned := map[string][]int{}

	for _, item := range set {
		m := item.(map[string]interface{})
		category := m["category"].(string)
		idsSet := m["sensor_ids"].(*schema.Set).List()
		ids := make([]int, len(idsSet))

		for i := range idsSet {
			ids[i] = idsSet[i].(int)
		}

		assigned[category] = ids
	}

	group.AssignedSensorIDs = &assigned

	return group
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(client.Client)

	group := groupFromResourceData(d, c)

	id, err := c.CreateGroup(ctx, group)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(id))

	return diags
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(client.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	group, err := c.ReadGroup(ctx, id)
	if err != nil {
		if _, ok := err.(client.ResourceNotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err := d.Set("name", group.Name); err != nil {
		return diag.FromErr(err)
	}

	if group.IsDefault != nil {
		if err := d.Set("is_default", *group.IsDefault); err != nil {
			return diag.FromErr(err)
		}
	}

	if group.PeriodicDailyReports != nil {
		if err := d.Set("periodic_daily_reports", *group.PeriodicDailyReports); err != nil {
			return diag.FromErr(err)
		}
	}

	if group.PeriodicWeeklyReports != nil {
		if err := d.Set("periodic_weekly_reports", *group.PeriodicWeeklyReports); err != nil {
			return diag.FromErr(err)
		}
	}

	if group.PeriodicMonthlyReports != nil {
		if err := d.Set("periodic_monthly_reports", *group.PeriodicMonthlyReports); err != nil {
			return diag.FromErr(err)
		}
	}

	if group.ArchivedServicesInPeriodicReports != nil {
		if err := d.Set("archived_services_in_periodic_reports", *group.ArchivedServicesInPeriodicReports); err != nil {
			return diag.FromErr(err)
		}
	}

	if group.AssignedSensorIDs != nil {
		list := make([]map[string]interface{}, 0, len(*group.AssignedSensorIDs))
		for category, ids := range *group.AssignedSensorIDs {
			list = append(list, map[string]interface{}{
				"category":   category,
				"sensor_ids": ids,
			})
		}

		if err := d.Set("assigned_sensor_ids", list); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(client.Client)

	group := groupFromResourceData(d, c)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.UpdateGroup(ctx, id, group)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(client.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteGroup(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
