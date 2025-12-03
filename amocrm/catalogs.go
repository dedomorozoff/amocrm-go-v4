package amocrm

import "context"

// Catalog represents an AmoCRM catalog
type Catalog struct {
	ID              int    `json:"id,omitempty"`
	Name            string `json:"name"`
	CreatedBy       int    `json:"created_by,omitempty"`
	UpdatedBy       int    `json:"updated_by,omitempty"`
	CreatedAt       int64  `json:"created_at,omitempty"`
	UpdatedAt       int64  `json:"updated_at,omitempty"`
	Sort            int    `json:"sort,omitempty"`
	Type            string `json:"type,omitempty"`
	CanAddElements  bool   `json:"can_add_elements,omitempty"`
	CanShowInCards  bool   `json:"can_show_in_cards,omitempty"`
	CanLinkMultiple bool   `json:"can_link_multiple,omitempty"`
	CanBeDeleted    bool   `json:"can_be_deleted,omitempty"`
	SDKWidgetCode   string `json:"sdk_widget_code,omitempty"`
	AccountID       int    `json:"account_id,omitempty"`
}

// CatalogsService handles communication with catalog-related methods
type CatalogsService struct {
	client *Client
}

// CatalogsResponse represents the API response for catalogs list
type CatalogsResponse struct {
	Embedded struct {
		Catalogs []Catalog `json:"catalogs"`
	} `json:"_embedded"`
	Links Links `json:"_links"`
}

// List retrieves a list of catalogs
func (s *CatalogsService) List(ctx context.Context) ([]Catalog, error) {
	var resp CatalogsResponse
	if err := s.client.GetJSON(ctx, "/catalogs", &resp); err != nil {
		return nil, err
	}

	return resp.Embedded.Catalogs, nil
}
