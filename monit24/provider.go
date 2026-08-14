package monit24

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/monit24/terraform-provider-monit24/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"user": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("MONIT24_USER", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("MONIT24_PASSWORD", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"monit24_group":                   resourceGroup(),
			"monit24_group_share":             resourceGroupShare(),
			"monit24_notification_address":    resourceNotificationAddress(),
			"monit24_periodic_report_address": resourcePeriodicReportAddress(),
			"monit24_service":                 resourceService(),
			"monit24_suspension":              resourceSuspension(),
			"monit24_weekly_suspension":       resourceWeeklySuspension(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	user := d.Get("user").(string)
	password := d.Get("password").(string)

	var diags diag.Diagnostics

	if user != "" && password != "" {
		c, err := client.NewBasicAuthClient(ctx, user, password)
		if err != nil {
			return nil, diag.FromErr(err)
		}

		return c, diags
	}

	return client.Client{}, diag.FromErr(errors.New("no credentials provided"))
}
