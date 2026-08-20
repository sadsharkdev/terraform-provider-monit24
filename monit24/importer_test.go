package monit24

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/monit24/terraform-provider-monit24/client"
)

// TestGroupImportRoundTrip exercises a simple numeric-ID resource end to
// end: ImportStatePassthroughContext just carries the id through, but the
// subsequent Read (which every real `terraform import` triggers) must
// still succeed against that id.
func TestGroupImportRoundTrip(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/groups/1" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(client.Group{Name: "imported group"})
	}))
	defer server.Close()

	c := client.NewTestClient(server.URL)

	d := schema.TestResourceDataRaw(t, resourceGroup().Schema, map[string]interface{}{})
	d.SetId("1")

	results, err := resourceGroup().Importer.StateContext(context.Background(), d, c)
	if err != nil {
		t.Fatalf("unexpected error from Importer.StateContext: %v", err)
	}
	if len(results) != 1 || results[0].Id() != "1" {
		t.Fatalf("expected passthrough import to return one ResourceData with id \"1\", got %+v", results)
	}

	if diags := resourceGroupRead(context.Background(), results[0], c); diags.HasError() {
		t.Fatalf("unexpected error reading imported group: %v", diags)
	}
	if results[0].Get("name").(string) != "imported group" {
		t.Errorf("expected imported group's name to be populated by the post-import Read, got %q", results[0].Get("name"))
	}
}

// TestGroupShareImportRoundTrip exercises a composite-ID resource: the
// imported id string must survive passthrough and still parse correctly
// via parseGroupShareID in the subsequent Read.
func TestGroupShareImportRoundTrip(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/groups/1/shares/2" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(client.GroupShare{GroupID: 1, AccountID: 2})
	}))
	defer server.Close()

	c := client.NewTestClient(server.URL)

	d := schema.TestResourceDataRaw(t, resourceGroupShare().Schema, map[string]interface{}{})
	d.SetId("1:2")

	results, err := resourceGroupShare().Importer.StateContext(context.Background(), d, c)
	if err != nil {
		t.Fatalf("unexpected error from Importer.StateContext: %v", err)
	}
	if len(results) != 1 || results[0].Id() != "1:2" {
		t.Fatalf("expected passthrough import to return one ResourceData with id \"1:2\", got %+v", results)
	}

	if diags := resourceGroupShareRead(context.Background(), results[0], c); diags.HasError() {
		t.Fatalf("unexpected error reading imported group_share: %v", diags)
	}
	if results[0].Get("group_id").(int) != 1 || results[0].Get("account_id").(int) != 2 {
		t.Errorf("expected imported group_share's group_id/account_id to be populated by the post-import Read, got group_id=%v account_id=%v",
			results[0].Get("group_id"), results[0].Get("account_id"))
	}
}

// TestUserDataSettingImportRoundTrip covers the composite-ID edge case
// where the key portion itself contains a colon — parseUserDataSettingID
// only splits on the first one (see TestUserDataSettingIDKeyContainingColon).
func TestUserDataSettingImportRoundTrip(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/user_data/123456/settings/namespace:setting" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("dark")
	}))
	defer server.Close()

	c := client.NewTestClient(server.URL)

	d := schema.TestResourceDataRaw(t, resourceUserDataSetting().Schema, map[string]interface{}{})
	d.SetId("123456:namespace:setting")

	results, err := resourceUserDataSetting().Importer.StateContext(context.Background(), d, c)
	if err != nil {
		t.Fatalf("unexpected error from Importer.StateContext: %v", err)
	}
	if len(results) != 1 || results[0].Id() != "123456:namespace:setting" {
		t.Fatalf("expected passthrough import to return one ResourceData with id \"123456:namespace:setting\", got %+v", results)
	}

	if diags := resourceUserDataSettingRead(context.Background(), results[0], c); diags.HasError() {
		t.Fatalf("unexpected error reading imported user_data_setting: %v", diags)
	}
	if results[0].Get("key").(string) != "namespace:setting" {
		t.Errorf("expected key with embedded colon to survive import+read, got %q", results[0].Get("key"))
	}
	if results[0].Get("value").(string) != "dark" {
		t.Errorf("expected value to be populated by the post-import Read, got %q", results[0].Get("value"))
	}
}
