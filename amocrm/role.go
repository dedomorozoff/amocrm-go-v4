package amocrm

import (
	"context"
	"fmt"
)

// Role represents an AmoCRM user role
type Role struct {
	ID       int           `json:"id,omitempty"`
	Name     string        `json:"name"`
	Rights   *RoleRights   `json:"rights,omitempty"`
	Embedded *RoleEmbedded `json:"_embedded,omitempty"`
	Links    *Links        `json:"_links,omitempty"`
}

// RoleEmbedded represents embedded role data
type RoleEmbedded struct {
	Users []UserShort `json:"users,omitempty"`
}

// UserShort represents a short user info
type UserShort struct {
	ID int `json:"id"`
}

// RoleRights represents role rights
type RoleRights struct {
	Leads         *EntityRights  `json:"leads,omitempty"`
	Contacts      *EntityRights  `json:"contacts,omitempty"`
	Companies     *EntityRights  `json:"companies,omitempty"`
	Tasks         *TaskRights    `json:"tasks,omitempty"`
	MailAccess    bool           `json:"mail_access,omitempty"`
	CatalogAccess bool           `json:"catalog_access,omitempty"`
	StatusRights  []StatusRights `json:"status_rights,omitempty"`
}

// RoleService handles communication with role-related methods
type RoleService struct {
	client *Client
}

// RolesResponse represents the API response for roles list
type RolesResponse struct {
	Embedded struct {
		Roles []Role `json:"roles"`
	} `json:"_embedded"`
	Links      Links `json:"_links"`
	TotalItems int   `json:"_total_items"`
	Page       int   `json:"_page"`
	PageCount  int   `json:"_page_count"`
}

// RolesFilter represents filter parameters for roles list
type RolesFilter struct {
	With  string // comma-separated list: users
	Limit int
	Page  int
}

// List retrieves a list of roles
func (s *RoleService) List(ctx context.Context, filter *RolesFilter) ([]Role, error) {
	path := "/roles"

	if filter != nil {
		path += "?"
		if filter.With != "" {
			path += fmt.Sprintf("with=%s&", filter.With)
		}
		if filter.Limit > 0 {
			path += fmt.Sprintf("limit=%d&", filter.Limit)
		}
		if filter.Page > 0 {
			path += fmt.Sprintf("page=%d&", filter.Page)
		}
	}

	var resp RolesResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	return resp.Embedded.Roles, nil
}

// Get retrieves a role by ID
func (s *RoleService) Get(ctx context.Context, id int, with string) (*Role, error) {
	path := fmt.Sprintf("/roles/%d", id)
	if with != "" {
		path += fmt.Sprintf("?with=%s", with)
	}

	var role Role
	if err := s.client.GetJSON(ctx, path, &role); err != nil {
		return nil, err
	}

	return &role, nil
}

// Create creates a new role
func (s *RoleService) Create(ctx context.Context, role *Role) (*Role, error) {
	var resp RolesResponse
	if err := s.client.PostJSON(ctx, "/roles", role, &resp); err != nil {
		return nil, err
	}

	if len(resp.Embedded.Roles) == 0 {
		return nil, fmt.Errorf("no role returned from API")
	}

	return &resp.Embedded.Roles[0], nil
}

// CreateBatch creates multiple roles in one request
func (s *RoleService) CreateBatch(ctx context.Context, roles []*Role) ([]Role, error) {
	rolesValues := make([]Role, len(roles))
	for i, r := range roles {
		rolesValues[i] = *r
	}

	var resp RolesResponse
	if err := s.client.PostJSON(ctx, "/roles", rolesValues, &resp); err != nil {
		return nil, err
	}

	return resp.Embedded.Roles, nil
}

// Update updates an existing role
func (s *RoleService) Update(ctx context.Context, id int, role *Role) (*Role, error) {
	path := fmt.Sprintf("/roles/%d", id)

	var updatedRole Role
	if err := s.client.PatchJSON(ctx, path, role, &updatedRole); err != nil {
		return nil, err
	}

	return &updatedRole, nil
}

// Delete deletes a role by ID
func (s *RoleService) Delete(ctx context.Context, id int) error {
	path := fmt.Sprintf("/roles/%d", id)
	return s.client.DeleteJSON(ctx, path)
}
