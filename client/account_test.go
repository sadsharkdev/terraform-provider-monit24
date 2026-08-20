package client

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

// SubaccountCreateRequest embeds Account; this locks in that Go flattens
// the embedded struct's fields into the same top-level JSON object rather
// than nesting them under an "Account" key.
func TestSubaccountCreateRequestFlattensEmbeddedAccount(t *testing.T) {
	req := SubaccountCreateRequest{
		Account: Account{
			Name:     "Client A",
			Username: "client-a",
		},
		UserData: UserData{
			EmailAddress: "client-a@example.com",
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if raw["name"] != "Client A" {
		t.Errorf("expected top-level \"name\" field, got %v", raw["name"])
	}
	if raw["username"] != "client-a" {
		t.Errorf("expected top-level \"username\" field, got %v", raw["username"])
	}
	if _, present := raw["Account"]; present {
		t.Error("expected Account to be flattened, not nested under an \"Account\" key")
	}

	userData, ok := raw["user_data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected user_data to be a nested object, got %v", raw["user_data"])
	}
	if userData["email_address"] != "client-a@example.com" {
		t.Errorf("expected user_data.email_address, got %v", userData["email_address"])
	}
}

func TestCreateSubaccountHitsSubaccountPath(t *testing.T) {
	var gotMethod, gotPath string

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 1}`))
	})

	id, err := c.CreateSubaccount(context.Background(), SubaccountCreateRequest{
		Account:  Account{Name: "Client A", Username: "client-a"},
		UserData: UserData{EmailAddress: "client-a@example.com"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Errorf("expected method POST, got %s", gotMethod)
	}
	if gotPath != "/accounts/subaccount" {
		t.Errorf("expected path /accounts/subaccount, got %s", gotPath)
	}
	if id != 1 {
		t.Errorf("expected id 1, got %d", id)
	}
}

func TestCreateAccountUserHitsAccountsPath(t *testing.T) {
	var gotPath string

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 2}`))
	})

	id, err := c.CreateAccountUser(context.Background(), SubaccountCreateRequest{
		Account:  Account{Name: "Colleague", Username: "colleague"},
		UserData: UserData{EmailAddress: "colleague@example.com"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotPath != "/accounts" {
		t.Errorf("expected path /accounts (not /accounts/subaccount), got %s", gotPath)
	}
	if id != 2 {
		t.Errorf("expected id 2, got %d", id)
	}
}

func TestReadAccount(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/accounts/1" {
			t.Errorf("expected path /accounts/1, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"name": "Client A",
			"username": "client-a",
			"package_id": 5,
			"is_activated": true,
			"is_blocked": false
		}`))
	})

	account, err := c.ReadAccount(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if account.Name != "Client A" || account.Username != "client-a" {
		t.Errorf("response didn't unmarshal correctly: %+v", account)
	}
	if account.PackageID == nil || *account.PackageID != 5 {
		t.Errorf("expected package_id=5, got %v", account.PackageID)
	}
	if account.IsActivated == nil || !*account.IsActivated {
		t.Errorf("expected is_activated=true, got %v", account.IsActivated)
	}
}

func TestReadAccountMissingPackageIDStaysNil(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name": "Client A", "username": "client-a"}`))
	})

	account, err := c.ReadAccount(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if account.PackageID != nil {
		t.Errorf("expected PackageID to stay nil when omitted from the response, got %v", *account.PackageID)
	}
}

func TestReadAccountNotFound(t *testing.T) {
	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{}`))
	})

	_, err := c.ReadAccount(context.Background(), 999)
	if _, ok := err.(ResourceNotFound); !ok {
		t.Fatalf("expected ResourceNotFound, got %T: %v", err, err)
	}
}

func TestUpdateAccount(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody Account

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.UpdateAccount(context.Background(), 1, Account{Name: "Renamed"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotMethod != http.MethodPut {
		t.Errorf("expected method PUT, got %s", gotMethod)
	}
	if gotPath != "/accounts/1" {
		t.Errorf("expected path /accounts/1, got %s", gotPath)
	}
	if gotBody.Name != "Renamed" {
		t.Errorf("expected request body name=Renamed, got %q", gotBody.Name)
	}
}

func TestChangeAccountPassword(t *testing.T) {
	var gotMethod, gotPath string
	var gotBody passwordUpdateData

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.ChangeAccountPassword(context.Background(), 1, "new-secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Errorf("expected method POST, got %s", gotMethod)
	}
	if gotPath != "/accounts/1/change_password" {
		t.Errorf("expected path /accounts/1/change_password, got %s", gotPath)
	}
	if gotBody.NewPassword != "new-secret" {
		t.Errorf("expected new_password=new-secret in request body, got %q", gotBody.NewPassword)
	}
}

func TestDeleteAccount(t *testing.T) {
	var gotMethod, gotPath string

	c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	})

	err := c.DeleteAccount(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotMethod != http.MethodDelete {
		t.Errorf("expected method DELETE, got %s", gotMethod)
	}
	if gotPath != "/accounts/1" {
		t.Errorf("expected path /accounts/1, got %s", gotPath)
	}
}
