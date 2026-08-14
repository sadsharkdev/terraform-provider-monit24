package monit24

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/monit24/terraform-provider-monit24/client"
)

func resourceService() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceCreate,
		ReadContext:   resourceServiceRead,
		UpdateContext: resourceServiceUpdate,
		DeleteContext: resourceServiceDelete,
		Schema: map[string]*schema.Schema{
			"type_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  600,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"is_active": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"is_archived": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"sensor_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"step_names": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"notification_channel_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"notification_condition_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"notification_mode_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
			"recovery_notification_mode_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "default",
			},
			"extended_settings": {
				Type:     schema.TypeMap,
				Computed: true,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func newServiceFromResourceData(service client.Service, d *schema.ResourceData) (client.Service, error) {
	service.Name = d.Get("name").(string)
	service.Address = d.Get("address").(string)
	service.GroupID = d.Get("group_id").(int)
	service.Interval = d.Get("interval").(int)

	if v, ok := d.GetOk("description"); ok {
		service.Description = strPtr(v.(string))
	}

	if v, ok := d.GetOk("interval"); ok {
		service.Interval = v.(int)
	}

	isActive := d.Get("is_active")
	service.IsActive = boolPtr(isActive.(bool))

	service.IsArchived = boolPtr(d.Get("is_archived").(bool))

	if v, ok := d.GetOk("sensor_ids"); ok {
		list := v.(*schema.Set).List()
		ids := make([]int, len(list))

		for i := range list {
			ids[i] = list[i].(int)
		}

		service.SensorIDs = &ids
	}

	if v, ok := d.GetOk("step_names"); ok {
		list := v.([]interface{})
		names := make([]string, len(list))

		for i := range list {
			names[i] = list[i].(string)
		}

		service.StepNames = &names
	}

	if v, ok := d.GetOk("notification_channel_ids"); ok {
		list := v.(*schema.Set).List()
		ids := make([]string, len(list))

		for i := range list {
			ids[i] = list[i].(string)
		}

		service.NotificationChannelIDs = &ids
	}

	if v, ok := d.GetOk("notification_condition_ids"); ok {
		list := v.(*schema.Set).List()
		ids := make([]string, len(list))

		for i := range list {
			ids[i] = list[i].(string)
		}

		service.NotificationConditionIDs = &ids
	}

	if v, ok := d.GetOk("notification_mode_id"); ok {
		service.NotificationModeID = strPtr(v.(string))
	}

	if v, ok := d.GetOk("recovery_notification_mode_id"); ok {
		service.RecoveryNotificationModeID = strPtr(v.(string))
	}

	if v, ok := d.GetOk("extended_settings"); ok {
		if service.ExtendedSettings == nil {
			service.ExtendedSettings = &map[string]interface{}{}
		}
		settings := v.(map[string]interface{})

		for k, v := range settings {
			stringValue := v.(string)

			intValue, err := strconv.Atoi(stringValue)
			if err == nil {
				(*service.ExtendedSettings)[k] = intValue
				continue
			}

			boolValue, err := parseBool(stringValue)
			if err == nil {
				(*service.ExtendedSettings)[k] = boolValue
				continue
			}

			(*service.ExtendedSettings)[k] = stringValue
		}
	}

	return service, nil
}

func resourceServiceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Client)

	s := client.Service{
		TypeID:  d.Get("type_id").(string),
		OwnerID: c.OwnerID(),
	}

	service, err := newServiceFromResourceData(s, d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := c.CreateService(ctx, service)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(id))

	return resourceServiceUpdate(ctx, d, m)
}

func resourceServiceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(client.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	service, err := c.ReadService(ctx, id)
	if err != nil {
		if _, ok := err.(client.ResourceNotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err := d.Set("type_id", service.TypeID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", service.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("address", service.Address); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("group_id", service.GroupID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("interval", service.Interval); err != nil {
		return diag.FromErr(err)
	}

	if service.Description != nil {
		if err := d.Set("description", *service.Description); err != nil {
			return diag.FromErr(err)
		}
	}

	if service.IsActive != nil {
		if err := d.Set("is_active", *service.IsActive); err != nil {
			return diag.FromErr(err)
		}
	}

	if service.IsArchived != nil {
		if err := d.Set("is_archived", *service.IsArchived); err != nil {
			return diag.FromErr(err)
		}
	}

	if service.SensorIDs != nil {
		if err := d.Set("sensor_ids", *service.SensorIDs); err != nil {
			return diag.FromErr(err)
		}
	}

	if service.StepNames != nil {
		if err := d.Set("step_names", *service.StepNames); err != nil {
			return diag.FromErr(err)
		}
	}

	if service.NotificationChannelIDs != nil {
		if err := d.Set("notification_channel_ids", *service.NotificationChannelIDs); err != nil {
			return diag.FromErr(err)
		}
	}

	if service.NotificationConditionIDs != nil {
		if err := d.Set("notification_condition_ids", *service.NotificationConditionIDs); err != nil {
			return diag.FromErr(err)
		}
	}

	if service.NotificationModeID != nil {
		if err := d.Set("notification_mode_id", *service.NotificationModeID); err != nil {
			return diag.FromErr(err)
		}
	}

	if service.RecoveryNotificationModeID != nil {
		if err := d.Set("recovery_notification_mode_id", *service.RecoveryNotificationModeID); err != nil {
			return diag.FromErr(err)
		}
	}

	if service.ExtendedSettings != nil {
		v, ok := d.GetOk("extended_settings")
		var definedSettings map[string]interface{}
		if ok {
			s := v.(map[string]interface{})
			if s != nil {
				definedSettings = s
			}
		}
		settings := mergeMaps(*service.ExtendedSettings, definedSettings)
		if err := d.Set("extended_settings", settings); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}

func resourceServiceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	s, err := c.ReadService(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	service, err := newServiceFromResourceData(s, d)
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.UpdateService(ctx, id, service)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceServiceRead(ctx, d, m)
}

func resourceServiceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	c := m.(client.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteService(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func parseBool(v string) (bool, error) {
	if v == "true" {
		return true, nil
	} else if v == "false" {
		return false, nil
	} else {
		return false, fmt.Errorf("failed to parse bool: %v", v)
	}
}

func mergeMaps(externalSettings, definedSettings map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{}

	for k, v := range externalSettings {
		_, ok := definedSettings[k]
		if !ok {
			continue
		}
		if v != nil {
			result[k] = fmt.Sprintf("%v", v)
		}
	}

	return result
}

func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
