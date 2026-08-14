package client

import (
	"context"
	"encoding/json"
	"fmt"
)

// The API stores settings as arbitrary JSON (string, number, bool, object,
// array). This client only ever writes plain JSON strings, so round-tripping
// a value this resource itself wrote is safe — reading a setting some other
// client wrote as a non-string JSON value will fail to unmarshal.
func (c Client) ReadUserDataSetting(ctx context.Context, id int, key string) (string, error) {
	resp, err := c.get(ctx, fmt.Sprintf("/user_data/%v/settings/%v", id, key))
	if err != nil {
		return "", err
	}

	var value string
	err = json.Unmarshal(resp, &value)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (c Client) PutUserDataSetting(ctx context.Context, id int, key string, value string) error {
	return c.put(ctx, fmt.Sprintf("/user_data/%v/settings/%v", id, key), value)
}

func (c Client) DeleteUserDataSetting(ctx context.Context, id int, key string) error {
	return c.delete(ctx, fmt.Sprintf("/user_data/%v/settings/%v", id, key))
}
