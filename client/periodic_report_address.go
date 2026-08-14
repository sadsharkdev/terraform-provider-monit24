package client

import (
	"context"
	"encoding/json"
	"fmt"
)

type PeriodicReportAddress struct {
	Address         string `json:"address"`
	ReportFrequency string `json:"report_frequency"`
	GroupID         int    `json:"group_id,omitempty"`
	OwnerID         int    `json:"owner_id"`
}

type CreatePeriodicReportAddressResponse struct {
	ID int `json:"id"`
}

func (c Client) CreatePeriodicReportAddress(ctx context.Context, req PeriodicReportAddress) (int, error) {
	resp, err := c.post(ctx, "/periodic_report_addresses", req)
	if err != nil {
		return 0, err
	}

	var response CreatePeriodicReportAddressResponse
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return 0, err
	}

	return response.ID, err
}

func (c Client) ReadPeriodicReportAddress(ctx context.Context, id int) (PeriodicReportAddress, error) {
	resp, err := c.get(ctx, fmt.Sprintf("/periodic_report_addresses/%v", id))
	if err != nil {
		return PeriodicReportAddress{}, err
	}

	var address PeriodicReportAddress
	err = json.Unmarshal(resp, &address)
	if err != nil {
		return PeriodicReportAddress{}, err
	}

	return address, nil
}

func (c Client) UpdatePeriodicReportAddress(ctx context.Context, id int, req PeriodicReportAddress) error {
	return c.put(ctx, fmt.Sprintf("/periodic_report_addresses/%v", id), req)
}

func (c Client) DeletePeriodicReportAddress(ctx context.Context, id int) error {
	return c.delete(ctx, fmt.Sprintf("/periodic_report_addresses/%v", id))
}
