package amocrm

import (
	"context"
	"fmt"
)

// User represents an AmoCRM user
type User struct {
	ID          int           `json:"id,omitempty"`
	Name        string        `json:"name"`
	Email       string        `json:"email"`
	Password    string        `json:"password,omitempty"`
	Lang        string        `json:"lang,omitempty"`
	Rights      *Rights       `json:"rights,omitempty"`
	UUID        *string       `json:"uuid,omitempty"`
	AmojoID     *string       `json:"amojo_id,omitempty"`
	UserRank    *string       `json:"user_rank,omitempty"`
	PhoneNumber *string       `json:"phone_number,omitempty"`
	Embedded    *UserEmbedded `json:"_embedded,omitempty"`
	Links       *Links        `json:"_links,omitempty"`
}

// UserEmbedded represents embedded user data
type UserEmbedded struct {
	Roles  []RoleShort `json:"roles,omitempty"`
	Groups []Group     `json:"groups,omitempty"`
}

// RoleShort represents a short role info
type RoleShort struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Links *Links `json:"_links,omitempty"`
}

// Rights represents user rights
type Rights struct {
	Leads         *EntityRights  `json:"leads,omitempty"`
	Contacts      *EntityRights  `json:"contacts,omitempty"`
	Companies     *EntityRights  `json:"companies,omitempty"`
	Tasks         *TaskRights    `json:"tasks,omitempty"`
	MailAccess    bool           `json:"mail_access,omitempty"`
	CatalogAccess bool           `json:"catalog_access,omitempty"`
	StatusRights  []StatusRights `json:"status_rights,omitempty"`
	IsAdmin       bool           `json:"is_admin,omitempty"`
	IsFree        bool           `json:"is_free,omitempty"`
	IsActive      bool           `json:"is_active,omitempty"`
	GroupID       *int           `json:"group_id,omitempty"`
	RoleID        *int           `json:"role_id,omitempty"`
}

// EntityRights represents rights for entities (leads, contacts, companies)
type EntityRights struct {
	View   string `json:"view"`
	Edit   string `json:"edit"`
	Add    string `json:"add"`
	Delete string `json:"delete"`
	Export string `json:"export"`
}

// TaskRights represents rights for tasks
type TaskRights struct {
	Edit   string `json:"edit"`
	Delete string `json:"delete"`
}

// StatusRights represents rights for specific status
type StatusRights struct {
	EntityType string              `json:"entity_type"`
	PipelineID int                 `json:"pipeline_id"`
	StatusID   int                 `json:"status_id"`
	Rights     *StatusEntityRights `json:"rights"`
}

// StatusEntityRights represents entity rights for status
type StatusEntityRights struct {
	View   string `json:"view"`
	Edit   string `json:"edit"`
	Delete string `json:"delete"`
	Export string `json:"export,omitempty"`
}

// Group represents an AmoCRM user group
type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// UserService handles communication with account-related methods
type UserService struct {
	client *Client
}

// UsersResponse represents the API response for users list
type UsersResponse struct {
	Embedded struct {
		Users []User `json:"users"`
	} `json:"_embedded"`
	Links     Links `json:"_links"`
	Page      int   `json:"_page"`
	PageCount int   `json:"_page_count"`
}

type UsersFilter struct {
	With  string // comma-separated list: role, group, uuid, amojo_id, user_rank, phone_number
	Limit int
	Page  int
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(ctx context.Context, id int, with string) (*User, error) {
	path := fmt.Sprintf("/users/%d", id)
	if with != "" {
		path += fmt.Sprintf("?with=%s", with)
	}

	var user User
	if err := s.client.GetJSON(ctx, path, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// List retrieves a list of users
func (s *UserService) List(ctx context.Context, filter *UsersFilter) ([]User, error) {
	path := "/users"

	if filter != nil {
		path += "?"

		if filter.Limit > 0 {
			path += fmt.Sprintf("limit=%d&", filter.Limit)
		}
		if filter.Page > 0 {
			path += fmt.Sprintf("page=%d&", filter.Page)
		}
		if filter.With != "" {
			path += fmt.Sprintf("with=%s&", filter.With)
		}
	}

	var resp UsersResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	return resp.Embedded.Users, nil
}

// Create creates a new user
func (s *UserService) Create(ctx context.Context, user *User) (*User, error) {
	type request struct {
		Users []User `json:"users"`
	}

	req := request{
		Users: []User{*user},
	}

	var resp UsersResponse
	if err := s.client.PostJSON(ctx, "/users", req, &resp); err != nil {
		return nil, err
	}

	if len(resp.Embedded.Users) == 0 {
		return nil, fmt.Errorf("no user returned from API")
	}

	return &resp.Embedded.Users[0], nil
}

// CreateBatch creates multiple users in one request
func (s *UserService) CreateBatch(ctx context.Context, users []*User) ([]User, error) {
	type request struct {
		Users []User `json:"users"`
	}

	usersValues := make([]User, len(users))
	for i, t := range users {
		usersValues[i] = *t
	}

	req := request{
		Users: usersValues,
	}

	var resp UsersResponse
	if err := s.client.PostJSON(ctx, "/users", req, &resp); err != nil {
		return nil, err
	}

	return resp.Embedded.Users, nil
}
