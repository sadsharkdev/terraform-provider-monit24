package client

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c Client) ReadUserData(ctx context.Context, id int) (UserData, error) {
	resp, err := c.get(ctx, fmt.Sprintf("/user_data/%v", id))
	if err != nil {
		return UserData{}, err
	}

	var userData UserData
	err = json.Unmarshal(resp, &userData)
	if err != nil {
		return UserData{}, err
	}

	return userData, nil
}

func (c Client) UpdateUserData(ctx context.Context, id int, req UserData) error {
	return c.put(ctx, fmt.Sprintf("/user_data/%v", id), req)
}
