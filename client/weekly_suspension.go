package client

import (
	"context"
	"encoding/json"
	"fmt"
)

type MinuteOfWeek struct {
	DayOfWeek int `json:"day_of_week"`
	Hour      int `json:"hour"`
	Minute    int `json:"minute"`
}

type WeeklySuspension struct {
	ServiceID         int          `json:"service_id"`
	StartMinute       MinuteOfWeek `json:"start_minute"`
	EndMinute         MinuteOfWeek `json:"end_minute"`
	OnlyNotifications *bool        `json:"only_notifications,omitempty"`
	Description       *string      `json:"description,omitempty"`
}

type CreateWeeklySuspensionResponse struct {
	ID int `json:"id"`
}

func (c Client) CreateWeeklySuspension(ctx context.Context, req WeeklySuspension) (int, error) {
	resp, err := c.post(ctx, "/weekly_suspensions", req)
	if err != nil {
		return 0, err
	}

	var response CreateWeeklySuspensionResponse
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return 0, err
	}

	return response.ID, err
}

func (c Client) ReadWeeklySuspension(ctx context.Context, id int) (WeeklySuspension, error) {
	resp, err := c.get(ctx, fmt.Sprintf("/weekly_suspensions/%v", id))
	if err != nil {
		return WeeklySuspension{}, err
	}

	var suspension WeeklySuspension
	err = json.Unmarshal(resp, &suspension)
	if err != nil {
		return WeeklySuspension{}, err
	}

	return suspension, nil
}

func (c Client) UpdateWeeklySuspension(ctx context.Context, id int, req WeeklySuspension) error {
	return c.put(ctx, fmt.Sprintf("/weekly_suspensions/%v", id), req)
}

func (c Client) DeleteWeeklySuspension(ctx context.Context, id int) error {
	return c.delete(ctx, fmt.Sprintf("/weekly_suspensions/%v", id))
}
