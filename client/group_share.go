package client

import (
	"context"
	"encoding/json"
	"fmt"
)

type GroupShare struct {
	AccountID          int   `json:"account_id"`
	GroupID            int   `json:"group_id"`
	CanArchiveServices *bool `json:"can_archive_services,omitempty"`
	CanCreateServices  *bool `json:"can_create_services,omitempty"`
	CanDeleteServices  *bool `json:"can_delete_services,omitempty"`
	CanForceAnalyses   *bool `json:"can_force_analyses,omitempty"`
	CanModifyGroup     *bool `json:"can_modify_group,omitempty"`
	CanModifyServices  *bool `json:"can_modify_services,omitempty"`
}

func (c Client) PutGroupShare(ctx context.Context, groupID int, accountID int, req GroupShare) error {
	return c.put(ctx, fmt.Sprintf("/groups/%v/shares/%v", groupID, accountID), req)
}

func (c Client) ReadGroupShare(ctx context.Context, groupID int, accountID int) (GroupShare, error) {
	resp, err := c.get(ctx, fmt.Sprintf("/groups/%v/shares/%v", groupID, accountID))
	if err != nil {
		return GroupShare{}, err
	}

	var share GroupShare
	err = json.Unmarshal(resp, &share)
	if err != nil {
		return GroupShare{}, err
	}

	return share, nil
}

func (c Client) DeleteGroupShare(ctx context.Context, groupID int, accountID int) error {
	return c.delete(ctx, fmt.Sprintf("/groups/%v/shares/%v", groupID, accountID))
}
