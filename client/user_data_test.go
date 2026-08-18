package client

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestReadUserData(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/user_data/1" {
			t.Errorf("expected path /user_data/1, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"email_address": "client-a@example.com",
			"contact_person": "Jane Doe"
		}`))
	})

	userData, err := c.ReadUserData(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if userData.EmailAddress != "client-a@example.com" {
		t.Errorf("expected email_address, got %q", userData.EmailAddress)
	}
	if userData.ContactPerson == nil || *userData.ContactPerson != "Jane Doe" {
		t.Errorf("expected contact_person=Jane Doe, got %v", userData.ContactPerson)
	}
}

// Settings is json.RawMessage specifically so that neither shape breaks
// unmarshaling — this provider never reads the field's contents, but a
// type-mismatched settings value must not fail the whole ReadUserData call.
func TestReadUserDataSettingsFieldAcceptsEitherShape(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"settings as an object", `{"email_address": "a@example.com", "settings": {"theme": "dark"}}`},
		{"settings as an array", `{"email_address": "a@example.com", "settings": ["theme", "dark"]}`},
		{"settings omitted", `{"email_address": "a@example.com"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.body))
			})

			userData, err := c.ReadUserData(context.Background(), 1)
			if err != nil {
				t.Fatalf("unexpected error unmarshaling settings shape %q: %v", tt.name, err)
			}
			if userData.EmailAddress != "a@example.com" {
				t.Errorf("expected email_address to still unmarshal correctly, got %q", userData.EmailAddress)
			}
		})
	}
}

func TestReadUserDataNotFound(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{}`))
	})

	_, err := c.ReadUserData(context.Background(), 999)
	if _, ok := err.(ResourceNotFound); !ok {
		t.Fatalf("expected ResourceNotFound, got %T: %v", err, err)
	}
}

func TestUpdateUserData(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody UserData

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.UpdateUserData(context.Background(), 1, UserData{EmailAddress: "updated@example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotMethod != http.MethodPut {
		t.Errorf("expected method PUT, got %s", gotMethod)
	}
	if gotPath != "/user_data/1" {
		t.Errorf("expected path /user_data/1, got %s", gotPath)
	}
	if gotBody.EmailAddress != "updated@example.com" {
		t.Errorf("expected request body email_address=updated@example.com, got %q", gotBody.EmailAddress)
	}
}
