package client

import (
	"context"
	"encoding/json"
	"fmt"
)

type Suspension struct {
	ServiceID         int     `json:"service_id"`
	StartTime         *string `json:"start_time,omitempty"`
	EndTime           string  `json:"end_time"`
	OnlyNotifications *bool   `json:"only_notifications,omitempty"`
	Description       *string `json:"description,omitempty"`
}

type CreateSuspensionResponse struct {
	ID int `json:"id"`
}

func (c Client) CreateSuspension(ctx context.Context, req Suspension) (int, error) {
	resp, err := c.post(ctx, "/suspensions", req)
	if err != nil {
		return 0, err
	}

	var response CreateSuspensionResponse
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return 0, err
	}

	return response.ID, err
}

func (c Client) ReadSuspension(ctx context.Context, id int) (Suspension, error) {
	resp, err := c.get(ctx, fmt.Sprintf("/suspensions/%v", id))
	if err != nil {
		return Suspension{}, err
	}

	var suspension Suspension
	err = json.Unmarshal(resp, &suspension)
	if err != nil {
		return Suspension{}, err
	}

	return suspension, nil
}

func (c Client) UpdateSuspension(ctx context.Context, id int, req Suspension) error {
	return c.put(ctx, fmt.Sprintf("/suspensions/%v", id), req)
}

func (c Client) DeleteSuspension(ctx context.Context, id int) error {
	return c.delete(ctx, fmt.Sprintf("/suspensions/%v", id))
}
