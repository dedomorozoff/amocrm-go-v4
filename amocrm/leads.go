package amocrm

import (
	"context"
	"fmt"
)

// Lead represents an AmoCRM lead (deal)
type Lead struct {
	ID                 int                `json:"id,omitempty"`
	Name               string             `json:"name"`
	Price              int                `json:"price,omitempty"`
	ResponsibleUserID  int                `json:"responsible_user_id,omitempty"`
	GroupID            int                `json:"group_id,omitempty"`
	StatusID           int                `json:"status_id,omitempty"`
	PipelineID         int                `json:"pipeline_id,omitempty"`
	LossReasonID       int                `json:"loss_reason_id,omitempty"`
	CreatedBy          int                `json:"created_by,omitempty"`
	UpdatedBy          int                `json:"updated_by,omitempty"`
	CreatedAt          int64              `json:"created_at,omitempty"`
	UpdatedAt          int64              `json:"updated_at,omitempty"`
	ClosedAt           int64              `json:"closed_at,omitempty"`
	ClosestTaskAt      int64              `json:"closest_task_at,omitempty"`
	IsClosed           bool               `json:"is_closed,omitempty"`
	CustomFieldsValues []CustomFieldValue `json:"custom_fields_values,omitempty"`
	Score              int                `json:"score,omitempty"`
	AccountID          int                `json:"account_id,omitempty"`
	LaborCost          int                `json:"labor_cost,omitempty"`
	Links              *Links             `json:"_links,omitempty"`
	Embedded           *Embedded          `json:"_embedded,omitempty"`
}

// LeadsService handles communication with lead-related methods
type LeadsService struct {
	client *Client
}

// LeadsResponse represents the API response for leads list
type LeadsResponse struct {
	Embedded struct {
		Leads []Lead `json:"leads"`
	} `json:"_embedded"`
	Links Links `json:"_links"`
	Page  int   `json:"_page"`
}

// LeadsFilter represents filter options for listing leads
type LeadsFilter struct {
	Query      string
	Limit      int
	Page       int
	With       string // comma-separated list: contacts, catalog_elements, loss_reason
	Order      string // created_at, updated_at, id, closed_at
	StatusID   []int
	PipelineID int
	UpdatedAt  map[string]int64 // filter by updated_at: map["from"]=timestamp, map["to"]=timestamp
}

// ListWithResponse retrieves a list of leads with full response including pagination links
func (s *LeadsService) ListWithResponse(ctx context.Context, filter *LeadsFilter) (*LeadsResponse, error) {
	path := "/leads"

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
		if filter.PipelineID > 0 {
			path += fmt.Sprintf("filter[pipeline_id]=%d&", filter.PipelineID)
		}
		for _, statusID := range filter.StatusID {
			path += fmt.Sprintf("filter[statuses][0][status_id]=%d&", statusID)
		}
		if filter.UpdatedAt != nil {
			if from, ok := filter.UpdatedAt["from"]; ok {
				path += fmt.Sprintf("filter[updated_at][from]=%d&", from)
			}
			if to, ok := filter.UpdatedAt["to"]; ok {
				path += fmt.Sprintf("filter[updated_at][to]=%d&", to)
			}
		}
	}

	var resp LeadsResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// List retrieves a list of leads
func (s *LeadsService) List(ctx context.Context, filter *LeadsFilter) ([]Lead, error) {
	resp, err := s.ListWithResponse(ctx, filter)
	if err != nil {
		return nil, err
	}

	return resp.Embedded.Leads, nil
}

// GetByID retrieves a lead by ID
func (s *LeadsService) GetByID(ctx context.Context, id int) (*Lead, error) {
	path := fmt.Sprintf("/leads/%d", id)

	var lead Lead
	if err := s.client.GetJSON(ctx, path, &lead); err != nil {
		return nil, err
	}

	return &lead, nil
}

// Create creates a new lead
func (s *LeadsService) Create(ctx context.Context, lead *Lead) (*Lead, error) {
	type request struct {
		Leads []Lead `json:"leads"`
	}

	req := request{
		Leads: []Lead{*lead},
	}

	var resp LeadsResponse
	if err := s.client.PostJSON(ctx, "/leads", req, &resp); err != nil {
		return nil, err
	}

	if len(resp.Embedded.Leads) == 0 {
		return nil, fmt.Errorf("no lead returned from API")
	}

	return &resp.Embedded.Leads[0], nil
}

// CreateBatch creates multiple leads in one request
func (s *LeadsService) CreateBatch(ctx context.Context, leads []*Lead) ([]Lead, error) {
	type request struct {
		Leads []Lead `json:"leads"`
	}

	leadsValues := make([]Lead, len(leads))
	for i, l := range leads {
		leadsValues[i] = *l
	}

	req := request{
		Leads: leadsValues,
	}

	var resp LeadsResponse
	if err := s.client.PostJSON(ctx, "/leads", req, &resp); err != nil {
		return nil, err
	}

	return resp.Embedded.Leads, nil
}

// Update updates an existing lead
func (s *LeadsService) Update(ctx context.Context, lead *Lead) (*Lead, error) {
	if lead.ID == 0 {
		return nil, fmt.Errorf("lead ID is required for update")
	}

	type request struct {
		Leads []Lead `json:"leads"`
	}

	req := request{
		Leads: []Lead{*lead},
	}

	var resp LeadsResponse
	if err := s.client.PatchJSON(ctx, "/leads", req, &resp); err != nil {
		return nil, err
	}

	if len(resp.Embedded.Leads) == 0 {
		return nil, fmt.Errorf("no lead returned from API")
	}

	return &resp.Embedded.Leads[0], nil
}

// UpdateBatch updates multiple leads in one request
func (s *LeadsService) UpdateBatch(ctx context.Context, leads []*Lead) ([]Lead, error) {
	type request struct {
		Leads []Lead `json:"leads"`
	}

	leadsValues := make([]Lead, len(leads))
	for i, l := range leads {
		if l.ID == 0 {
			return nil, fmt.Errorf("lead ID is required for update at index %d", i)
		}
		leadsValues[i] = *l
	}

	req := request{
		Leads: leadsValues,
	}

	var resp LeadsResponse
	if err := s.client.PatchJSON(ctx, "/leads", req, &resp); err != nil {
		return nil, err
	}

	return resp.Embedded.Leads, nil
}

// LinkContacts links contacts to a lead
func (s *LeadsService) LinkContacts(ctx context.Context, leadID int, contactIDs []int) error {
	type linkRequest struct {
		ToEntityID   int    `json:"to_entity_id"`
		ToEntityType string `json:"to_entity_type"`
	}

	type request struct {
		Links []linkRequest `json:"links"`
	}

	links := make([]linkRequest, len(contactIDs))
	for i, contactID := range contactIDs {
		links[i] = linkRequest{
			ToEntityID:   contactID,
			ToEntityType: "contacts",
		}
	}

	req := request{Links: links}

	path := fmt.Sprintf("/leads/%d/link", leadID)
	return s.client.PostJSON(ctx, path, req, nil)
}

// LinkCompany links a company to a lead
func (s *LeadsService) LinkCompany(ctx context.Context, leadID int, companyID int) error {
	type linkRequest struct {
		ToEntityID   int    `json:"to_entity_id"`
		ToEntityType string `json:"to_entity_type"`
	}

	type request struct {
		Links []linkRequest `json:"links"`
	}

	req := request{
		Links: []linkRequest{
			{
				ToEntityID:   companyID,
				ToEntityType: "companies",
			},
		},
	}

	path := fmt.Sprintf("/leads/%d/link", leadID)
	return s.client.PostJSON(ctx, path, req, nil)
}
