package amocrm

import (
	"context"
	"fmt"
)

// Contact represents an AmoCRM contact
type Contact struct {
	ID                 int                `json:"id,omitempty"`
	Name               string             `json:"name"`
	FirstName          string             `json:"first_name,omitempty"`
	LastName           string             `json:"last_name,omitempty"`
	ResponsibleUserID  int                `json:"responsible_user_id,omitempty"`
	GroupID            int                `json:"group_id,omitempty"`
	CreatedBy          int                `json:"created_by,omitempty"`
	UpdatedBy          int                `json:"updated_by,omitempty"`
	CreatedAt          int64              `json:"created_at,omitempty"`
	UpdatedAt          int64              `json:"updated_at,omitempty"`
	ClosestTaskAt      int64              `json:"closest_task_at,omitempty"`
	CustomFieldsValues []CustomFieldValue `json:"custom_fields_values,omitempty"`
	AccountID          int                `json:"account_id,omitempty"`
	Links              *Links             `json:"_links,omitempty"`
	Embedded           *Embedded          `json:"_embedded,omitempty"`
}

// ContactsService handles communication with contact-related methods
type ContactsService struct {
	client *Client
}

// ContactsResponse represents the API response for contacts list
type ContactsResponse struct {
	Embedded struct {
		Contacts []Contact `json:"contacts"`
	} `json:"_embedded"`
	Links Links `json:"_links"`
	Page  Page  `json:"_page,omitempty"`
}

// Page represents pagination information
type Page struct {
	Size  int `json:"size,omitempty"`
	Count int `json:"count,omitempty"`
}

// ContactsFilter represents filter options for listing contacts
type ContactsFilter struct {
	Query string
	Limit int
	Page  int
	With  string // comma-separated list: leads, customers, catalog_elements
	Order string // created_at, updated_at, id
}

// List retrieves a list of contacts
func (s *ContactsService) List(ctx context.Context, filter *ContactsFilter) ([]Contact, error) {
	path := "/contacts"

	if filter != nil {
		path += "?"
		if filter.Query != "" {
			path += fmt.Sprintf("query=%s&", filter.Query)
		}
		if filter.Limit > 0 {
			path += fmt.Sprintf("limit=%d&", filter.Limit)
		}
		if filter.Page > 0 {
			path += fmt.Sprintf("page=%d&", filter.Page)
		}
		if filter.With != "" {
			path += fmt.Sprintf("with=%s&", filter.With)
		}
		if filter.Order != "" {
			path += fmt.Sprintf("order[%s]=asc&", filter.Order)
		}
	}

	var resp ContactsResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	return resp.Embedded.Contacts, nil
}

// GetByID retrieves a contact by ID
func (s *ContactsService) GetByID(ctx context.Context, id int) (*Contact, error) {
	path := fmt.Sprintf("/contacts/%d", id)

	var contact Contact
	if err := s.client.GetJSON(ctx, path, &contact); err != nil {
		return nil, err
	}

	return &contact, nil
}

// Create creates a new contact
func (s *ContactsService) Create(ctx context.Context, contact *Contact) (*Contact, error) {
	type request struct {
		Contacts []Contact `json:"contacts"`
	}

	req := request{
		Contacts: []Contact{*contact},
	}

	var resp ContactsResponse
	if err := s.client.PostJSON(ctx, "/contacts", req, &resp); err != nil {
		return nil, err
	}

	if len(resp.Embedded.Contacts) == 0 {
		return nil, fmt.Errorf("no contact returned from API")
	}

	return &resp.Embedded.Contacts[0], nil
}

// CreateBatch creates multiple contacts in one request
func (s *ContactsService) CreateBatch(ctx context.Context, contacts []*Contact) ([]Contact, error) {
	type request struct {
		Contacts []Contact `json:"contacts"`
	}

	// Convert pointers to values
	contactsValues := make([]Contact, len(contacts))
	for i, c := range contacts {
		contactsValues[i] = *c
	}

	req := request{
		Contacts: contactsValues,
	}

	var resp ContactsResponse
	if err := s.client.PostJSON(ctx, "/contacts", req, &resp); err != nil {
		return nil, err
	}

	return resp.Embedded.Contacts, nil
}

// Update updates an existing contact
func (s *ContactsService) Update(ctx context.Context, contact *Contact) (*Contact, error) {
	if contact.ID == 0 {
		return nil, fmt.Errorf("contact ID is required for update")
	}

	type request struct {
		Contacts []Contact `json:"contacts"`
	}

	req := request{
		Contacts: []Contact{*contact},
	}

	var resp ContactsResponse
	if err := s.client.PatchJSON(ctx, "/contacts", req, &resp); err != nil {
		return nil, err
	}

	if len(resp.Embedded.Contacts) == 0 {
		return nil, fmt.Errorf("no contact returned from API")
	}

	return &resp.Embedded.Contacts[0], nil
}

// UpdateBatch updates multiple contacts in one request
func (s *ContactsService) UpdateBatch(ctx context.Context, contacts []*Contact) ([]Contact, error) {
	type request struct {
		Contacts []Contact `json:"contacts"`
	}

	// Convert pointers to values
	contactsValues := make([]Contact, len(contacts))
	for i, c := range contacts {
		if c.ID == 0 {
			return nil, fmt.Errorf("contact ID is required for update at index %d", i)
		}
		contactsValues[i] = *c
	}

	req := request{
		Contacts: contactsValues,
	}

	var resp ContactsResponse
	if err := s.client.PatchJSON(ctx, "/contacts", req, &resp); err != nil {
		return nil, err
	}

	return resp.Embedded.Contacts, nil
}
