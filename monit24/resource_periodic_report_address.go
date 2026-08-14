package monit24

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/monit24/terraform-provider-monit24/client"
)

func resourcePeriodicReportAddress() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePeriodicReportAddressCreate,
		ReadContext:   resourcePeriodicReportAddressRead,
		UpdateContext: resourcePeriodicReportAddressUpdate,
		DeleteContext: resourcePeriodicReportAddressDelete,
		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"report_frequency": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func periodicReportAddressFromResourceData(d *schema.ResourceData, c client.Client) client.PeriodicReportAddress {
	return client.PeriodicReportAddress{
		Address:         d.Get("address").(string),
		ReportFrequency: d.Get("report_frequency").(string),
		GroupID:         d.Get("group_id").(int),
		OwnerID:         c.OwnerID(),
	}
}

func resourcePeriodicReportAddressCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(client.Client)

	address := periodicReportAddressFromResourceData(d, c)

	id, err := c.CreatePeriodicReportAddress(ctx, address)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(id))

	return diags
}

func resourcePeriodicReportAddressRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(client.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	address, err := c.ReadPeriodicReportAddress(ctx, id)
	if err != nil {
		if _, ok := err.(client.ResourceNotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err := d.Set("address", address.Address); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("report_frequency", address.ReportFrequency); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("group_id", address.GroupID); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourcePeriodicReportAddressUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(client.Client)

	address := periodicReportAddressFromResourceData(d, c)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.UpdatePeriodicReportAddress(ctx, id, address)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourcePeriodicReportAddressDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(client.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeletePeriodicReportAddress(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
