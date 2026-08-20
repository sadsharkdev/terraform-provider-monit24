package client

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestPutGroupShare(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody GroupShare

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusNoContent)
	})

	canModifyGroup := true

	err := c.PutGroupShare(context.Background(), 1, 2, GroupShare{
		GroupID:        1,
		AccountID:      2,
		CanModifyGroup: &canModifyGroup,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodPut || gotPath != "/groups/1/shares/2" {
		t.Errorf("expected PUT /groups/1/shares/2, got %s %s", gotMethod, gotPath)
	}
	if gotBody.CanModifyGroup == nil || !*gotBody.CanModifyGroup {
		t.Errorf("expected can_modify_group=true in request body, got %v", gotBody.CanModifyGroup)
	}
}

func TestReadGroupShare(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/groups/1/shares/2" {
			t.Errorf("expected path /groups/1/shares/2, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"group_id": 1, "account_id": 2, "can_modify_group": true}`))
	})

	share, err := c.ReadGroupShare(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if share.GroupID != 1 || share.AccountID != 2 {
		t.Errorf("response didn't unmarshal correctly: %+v", share)
	}
}

func TestReadGroupShareNotFound(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{}`))
	})

	_, err := c.ReadGroupShare(context.Background(), 1, 999)
	if _, ok := err.(ResourceNotFound); !ok {
		t.Fatalf("expected ResourceNotFound, got %T: %v", err, err)
	}
}

func TestDeleteGroupShare(t *testing.T) {
	var gotMethod, gotPath string

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.DeleteGroupShare(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotMethod != http.MethodDelete || gotPath != "/groups/1/shares/2" {
		t.Errorf("expected DELETE /groups/1/shares/2, got %s %s", gotMethod, gotPath)
	}
}
