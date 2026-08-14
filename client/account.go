package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserData struct {
	EmailAddress            string    `json:"email_address"`
	Address                 *string   `json:"address,omitempty"`
	ContactPerson           *string   `json:"contact_person,omitempty"`
	PhoneNumber             *string   `json:"phone_number,omitempty"`
	TaxIdentificationNumber *string   `json:"tax_identification_number,omitempty"`
	IPWhitelist             *[]string `json:"ip_whitelist,omitempty"`
	IPWhitelistEnabled      *bool     `json:"ip_whitelist_enabled,omitempty"`
	ID                      *int      `json:"id,omitempty"`
	CreatedAt               *string   `json:"created_at,omitempty"`
	Has2FAEnabled           *bool     `json:"has_2fa_enabled,omitempty"`
	Settings                *[]string `json:"settings,omitempty"`
}

type Account struct {
	Name                       string  `json:"name"`
	Username                   string  `json:"username"`
	PackageID                  int     `json:"package_id,omitempty"`
	ParentAccountID            *int    `json:"parent_account_id,omitempty"`
	IsActivated                *bool   `json:"is_activated,omitempty"`
	IsBlocked                  *bool   `json:"is_blocked,omitempty"`
	IsReadOnly                 *bool   `json:"is_read_only,omitempty"`
	DisableLegacyNotifications *bool   `json:"disable_legacy_notifications,omitempty"`
	LanguageID                 *string `json:"language_id,omitempty"`
	TimeZoneID                 *string `json:"time_zone_id,omitempty"`
}

type SubaccountCreateRequest struct {
	Account
	UserData           UserData `json:"user_data"`
	Password           *string  `json:"password,omitempty"`
	SetPasswordURL     *string  `json:"set_password_url,omitempty"`
	SubaccountBlock    *bool    `json:"subaccount_block,omitempty"`
	SubaccountEdit     *bool    `json:"subaccount_edit,omitempty"`
	Is2FASetupRequired *bool    `json:"is_2fa_setup_required,omitempty"`
}

type CreateSubaccountResponse struct {
	ID int `json:"id"`
}

type passwordUpdateData struct {
	NewPassword string `json:"new_password"`
}

func (c Client) CreateSubaccount(ctx context.Context, req SubaccountCreateRequest) (int, error) {
	resp, err := c.post(ctx, "/accounts/subaccount", req)
	if err != nil {
		return 0, err
	}

	var response CreateSubaccountResponse
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return 0, err
	}

	return response.ID, err
}

func (c Client) CreateAccountUser(ctx context.Context, req SubaccountCreateRequest) (int, error) {
	resp, err := c.post(ctx, "/accounts", req)
	if err != nil {
		return 0, err
	}

	var response CreateSubaccountResponse
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return 0, err
	}

	return response.ID, err
}

func (c Client) ReadAccount(ctx context.Context, id int) (Account, error) {
	resp, err := c.get(ctx, fmt.Sprintf("/accounts/%v", id))
	if err != nil {
		return Account{}, err
	}

	var account Account
	err = json.Unmarshal(resp, &account)
	if err != nil {
		return Account{}, err
	}

	return account, nil
}

func (c Client) UpdateAccount(ctx context.Context, id int, req Account) error {
	return c.put(ctx, fmt.Sprintf("/accounts/%v", id), req)
}

func (c Client) ChangeAccountPassword(ctx context.Context, id int, newPassword string) error {
	_, err := c.requestWithPayload(ctx, "POST", fmt.Sprintf("/accounts/%v/change_password", id), passwordUpdateData{NewPassword: newPassword}, http.StatusNoContent)
	return err
}

func (c Client) DeleteAccount(ctx context.Context, id int) error {
	return c.delete(ctx, fmt.Sprintf("/accounts/%v", id))
}
