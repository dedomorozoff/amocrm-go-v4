package amocrm

import (
	"context"
	"fmt"
)

// Company represents an AmoCRM company
type Company struct {
	ID                 int                `json:"id,omitempty"`
	Name               string             `json:"name"`
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

// CompaniesService handles communication with company-related methods
type CompaniesService struct {
	client *Client
}

// CompaniesResponse represents the API response for companies list
type CompaniesResponse struct {
	Embedded struct {
		Companies []Company `json:"companies"`
	} `json:"_embedded"`
	Links Links `json:"_links"`
	Page  int   `json:"_page"`
}

// CompaniesFilter represents filter options for listing companies
type CompaniesFilter struct {
	Query string
	Limit int
	Page  int
	With  string // comma-separated list: leads, customers, contacts, catalog_elements
	Order string // created_at, updated_at, id
}

// List retrieves a list of companies
func (s *CompaniesService) List(ctx context.Context, filter *CompaniesFilter) ([]Company, error) {
	path := "/companies"

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

	var resp CompaniesResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	return resp.Embedded.Companies, nil
}

// GetByID retrieves a company by ID
func (s *CompaniesService) GetByID(ctx context.Context, id int) (*Company, error) {
	path := fmt.Sprintf("/companies/%d", id)

	var company Company
	if err := s.client.GetJSON(ctx, path, &company); err != nil {
		return nil, err
	}

	return &company, nil
}

// Create creates a new company
func (s *CompaniesService) Create(ctx context.Context, company *Company) (*Company, error) {
	type request struct {
		Companies []Company `json:"companies"`
	}

	req := request{
		Companies: []Company{*company},
	}

	var resp CompaniesResponse
	if err := s.client.PostJSON(ctx, "/companies", req, &resp); err != nil {
		return nil, err
	}

	if len(resp.Embedded.Companies) == 0 {
		return nil, fmt.Errorf("no company returned from API")
	}

	return &resp.Embedded.Companies[0], nil
}

// CreateBatch creates multiple companies in one request
func (s *CompaniesService) CreateBatch(ctx context.Context, companies []*Company) ([]Company, error) {
	type request struct {
		Companies []Company `json:"companies"`
	}

	companiesValues := make([]Company, len(companies))
	for i, c := range companies {
		companiesValues[i] = *c
	}

	req := request{
		Companies: companiesValues,
	}

	var resp CompaniesResponse
	if err := s.client.PostJSON(ctx, "/companies", req, &resp); err != nil {
		return nil, err
	}

	return resp.Embedded.Companies, nil
}

// Update updates an existing company
func (s *CompaniesService) Update(ctx context.Context, company *Company) (*Company, error) {
	if company.ID == 0 {
		return nil, fmt.Errorf("company ID is required for update")
	}

	type request struct {
		Companies []Company `json:"companies"`
	}

	req := request{
		Companies: []Company{*company},
	}

	var resp CompaniesResponse
	if err := s.client.PatchJSON(ctx, "/companies", req, &resp); err != nil {
		return nil, err
	}

	if len(resp.Embedded.Companies) == 0 {
		return nil, fmt.Errorf("no company returned from API")
	}

	return &resp.Embedded.Companies[0], nil
}

// UpdateBatch updates multiple companies in one request
func (s *CompaniesService) UpdateBatch(ctx context.Context, companies []*Company) ([]Company, error) {
	type request struct {
		Companies []Company `json:"companies"`
	}

	companiesValues := make([]Company, len(companies))
	for i, c := range companies {
		if c.ID == 0 {
			return nil, fmt.Errorf("company ID is required for update at index %d", i)
		}
		companiesValues[i] = *c
	}

	req := request{
		Companies: companiesValues,
	}

	var resp CompaniesResponse
	if err := s.client.PatchJSON(ctx, "/companies", req, &resp); err != nil {
		return nil, err
	}

	return resp.Embedded.Companies, nil
}
