package amocrm

import (
	"context"
	"fmt"
)

type Event struct {
	ID          int64                  `json:"id"`
	Type        string                 `json:"type"`
	EntityID    int                    `json:"entity_id"`
	EntityType  string                 `json:"entity_type"`
	CreatedBy   int                    `json:"created_by"`
	CreatedAt   int64                  `json:"created_at"`
	ValueBefore map[string]interface{} `json:"value_before"`
	ValueAfter  map[string]interface{} `json:"value_after"`
	AccountID   int                    `json:"account_id"`
}

type EventsService struct {
	client *Client
}

type EventsResponse struct {
	Embedded struct {
		Events []Event `json:"events"`
	} `json:"_embedded"`
	Links     Links `json:"_links"`
	Page      int   `json:"_page"`
	PageCount int   `json:"_page_count"`
}

type EventsFilter struct {
	Limit      int
	Page       int
	EntityType EntityType
	EntityID   int
	Type       []string
	CreatedAt  map[string]int64
}

func (s *EventsService) List(ctx context.Context, filter *EventsFilter) ([]Event, error) {
	path := "/events"

	if filter != nil {
		path += "?"
		if filter.Limit > 0 {
			path += fmt.Sprintf("limit=%d&", filter.Limit)
		}
		if filter.Page > 0 {
			path += fmt.Sprintf("page=%d&", filter.Page)
		}
		if filter.EntityType != "" {
			path += fmt.Sprintf("filter[entity]=%s&", filter.EntityType)
		}
		if filter.EntityID > 0 {
			path += fmt.Sprintf("filter[entity_id]=%d&", filter.EntityID)
		}
		for _, eventType := range filter.Type {
			path += fmt.Sprintf("filter[type][]=%s&", eventType)
		}
		if filter.CreatedAt != nil {
			if from, ok := filter.CreatedAt["from"]; ok {
				path += fmt.Sprintf("filter[created_at][from]=%d&", from)
			}
			if to, ok := filter.CreatedAt["to"]; ok {
				path += fmt.Sprintf("filter[created_at][to]=%d&", to)
			}
		}
	}

	var resp EventsResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	return resp.Embedded.Events, nil
}
