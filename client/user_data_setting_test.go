package client

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestReadUserDataSetting(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/user_data/1/settings/theme" {
			t.Errorf("expected path /user_data/1/settings/theme, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`"dark"`))
	})

	value, err := c.ReadUserDataSetting(context.Background(), 1, "theme")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if value != "dark" {
		t.Errorf("expected value=dark, got %q", value)
	}
}

func TestReadUserDataSettingNotFound(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{}`))
	})

	_, err := c.ReadUserDataSetting(context.Background(), 1, "missing")
	if _, ok := err.(ResourceNotFound); !ok {
		t.Fatalf("expected ResourceNotFound, got %T: %v", err, err)
	}
}

func TestPutUserDataSetting(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody string

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.PutUserDataSetting(context.Background(), 1, "theme", "dark")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotMethod != http.MethodPut {
		t.Errorf("expected method PUT, got %s", gotMethod)
	}
	if gotPath != "/user_data/1/settings/theme" {
		t.Errorf("expected path /user_data/1/settings/theme, got %s", gotPath)
	}
	if gotBody != "dark" {
		t.Errorf("expected request body \"dark\", got %q", gotBody)
	}
}

func TestDeleteUserDataSetting(t *testing.T) {
	var gotMethod, gotPath string

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.DeleteUserDataSetting(context.Background(), 1, "theme")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotMethod != http.MethodDelete {
		t.Errorf("expected method DELETE, got %s", gotMethod)
	}
	if gotPath != "/user_data/1/settings/theme" {
		t.Errorf("expected path /user_data/1/settings/theme, got %s", gotPath)
	}
}
