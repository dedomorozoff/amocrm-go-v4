package amocrm

import (
	"context"
)

// Account represents AmoCRM account information
type Account struct {
	ID                      int              `json:"id"`
	Name                    string           `json:"name"`
	Subdomain               string           `json:"subdomain"`
	CreatedAt               int64            `json:"created_at"`
	CreatedBy               int              `json:"created_by"`
	UpdatedAt               int64            `json:"updated_at"`
	UpdatedBy               int              `json:"updated_by"`
	CurrentUserID           int              `json:"current_user_id"`
	Country                 string           `json:"country"`
	Currency                string           `json:"currency"`
	CustomersMode           string           `json:"customers_mode"`
	IsUnsortedOn            bool             `json:"is_unsorted_on"`
	MobileFeatureVersion    int              `json:"mobile_feature_version"`
	IsLossReasonEnabled     bool             `json:"is_loss_reason_enabled"`
	IsHelpbotEnabled        bool             `json:"is_helpbot_enabled"`
	IsTechnicalAccount      bool             `json:"is_technical_account"`
	ContactNameDisplayOrder int              `json:"contact_name_display_order"`
	AmojoID                 string           `json:"amojo_id,omitempty"`
	UUID                    string           `json:"uuid,omitempty"`
	Version                 int              `json:"version,omitempty"`
	Embedded                *AccountEmbedded `json:"_embedded,omitempty"`
}

// AccountEmbedded represents embedded account data
type AccountEmbedded struct {
	Users  []User  `json:"users,omitempty"`
	Groups []Group `json:"groups,omitempty"`
}

// AccountService handles communication with account-related methods
type AccountService struct {
	client *Client
}

// Get retrieves account information
func (s *AccountService) Get(ctx context.Context) (*Account, error) {
	var account Account
	if err := s.client.GetJSON(ctx, "/account", &account); err != nil {
		return nil, err
	}
	return &account, nil
}

// GetWithUsers retrieves account information with users
func (s *AccountService) GetWithUsers(ctx context.Context) (*Account, error) {
	var account Account
	if err := s.client.GetJSON(ctx, "/account?with=users", &account); err != nil {
		return nil, err
	}
	return &account, nil
}

// GetWithUsersAndGroups retrieves account information with users and groups
func (s *AccountService) GetWithUsersAndGroups(ctx context.Context) (*Account, error) {
	var account Account
	if err := s.client.GetJSON(ctx, "/account?with=users,groups", &account); err != nil {
		return nil, err
	}
	return &account, nil
}
