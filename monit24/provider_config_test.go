package monit24

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestProviderInternalValidate(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("Provider().InternalValidate() failed: %v", err)
	}
}

func TestProviderResourcesRegistered(t *testing.T) {
	expected := []string{
		"monit24_account_user",
		"monit24_group",
		"monit24_group_share",
		"monit24_notification_address",
		"monit24_periodic_report_address",
		"monit24_service",
		"monit24_subaccount",
		"monit24_suspension",
		"monit24_user_data",
		"monit24_user_data_setting",
		"monit24_weekly_suspension",
	}

	resources := Provider().ResourcesMap

	if len(resources) != len(expected) {
		t.Errorf("expected %d registered resources, got %d: %v", len(expected), len(resources), resources)
	}

	for _, name := range expected {
		if resources[name] == nil {
			t.Errorf("expected resource %q to be registered in Provider().ResourcesMap", name)
		}
	}
}

func TestAllResourcesHaveImporter(t *testing.T) {
	for name, r := range Provider().ResourcesMap {
		if r.Importer == nil {
			t.Errorf("resource %q has no Importer configured — CLAUDE.md documents that all resources support import", name)
			continue
		}
		if r.Importer.StateContext == nil {
			t.Errorf("resource %q has an Importer but no StateContext function", name)
		}
	}
}

func TestProviderConfigureNoCredentialsReturnsError(t *testing.T) {
	d := schema.TestResourceDataRaw(t, Provider().Schema, map[string]interface{}{})

	_, diags := providerConfigure(context.Background(), d)

	if !diags.HasError() {
		t.Fatal("expected an error when no token and no user/password are configured, got none")
	}
}

func TestProviderConfigureIncompleteBasicAuthReturnsError(t *testing.T) {
	// user set without password (and vice versa) must not silently fall
	// through to an unauthenticated client — both are required together.
	d := schema.TestResourceDataRaw(t, Provider().Schema, map[string]interface{}{
		"user": "someone",
	})

	_, diags := providerConfigure(context.Background(), d)

	if !diags.HasError() {
		t.Fatal("expected an error when user is set without password, got none")
	}
}
