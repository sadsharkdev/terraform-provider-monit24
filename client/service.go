package client

import (
	"context"
	"encoding/json"
	"fmt"
)

type Service struct {
	TypeID                     string                  `json:"type_id"`
	Name                       string                  `json:"name"`
	Address                    string                  `json:"address"`
	GroupID                    int                     `json:"group_id,omitempty"`
	Interval                   int                     `json:"interval"`
	Description                *string                 `json:"description,omitempty"`
	IsActive                   *bool                   `json:"is_active,omitempty"`
	SensorIDs                  *[]int                  `json:"sensor_ids,omitempty"`
	StepNames                  *[]string               `json:"step_names,omitempty"`
	NotificationChannelIDs     *[]string               `json:"notification_channel_ids,omitempty"`
	NotificationConditionIDs   *[]string               `json:"notification_condition_ids,omitempty"`
	NotificationModeID         *string                 `json:"notification_mode_id,omitempty"`
	RecoveryNotificationModeID *string                 `json:"recovery_notification_mode_id,omitempty"`
	ExtendedSettings           *map[string]interface{} `json:"extended_settings,omitempty"`
	OwnerID                    int                     `json:"owner_id"`
}

type CreateServiceResponse struct {
	ID int `json:"id"`
}

func (c Client) CreateService(ctx context.Context, req Service) (int, error) {
	resp, err := c.post(ctx, "/services", req)
	if err != nil {
		return 0, err
	}

	var response CreateServiceResponse
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return 0, err
	}

	return response.ID, err
}

func (c Client) ReadService(ctx context.Context, id int) (Service, error) {
	resp, err := c.get(ctx, fmt.Sprintf("/services/%v", id))
	if err != nil {
		return Service{}, err
	}

	var service Service
	err = json.Unmarshal(resp, &service)
	if err != nil {
		return Service{}, err
	}

	return service, nil
}

func (c Client) UpdateService(ctx context.Context, id int, req Service) error {
	return c.put(ctx, fmt.Sprintf("/services/%v", id), req)
}

func (c Client) DeleteService(ctx context.Context, id int) error {
	return c.delete(ctx, fmt.Sprintf("/services/%v", id))
}
