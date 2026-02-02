package amocrm

import (
	"context"
	"fmt"
)

type TagsService struct {
	client *Client
}

type TagsResponse struct {
	Embedded struct {
		Tags []Tag `json:"tags"`
	} `json:"_embedded"`
	Links Links `json:"_links"`
	Page  int   `json:"_page"`
}

type TagsFilter struct {
	Limit int
	Page  int
}

func (s *TagsService) List(ctx context.Context, entityType EntityType, filter *TagsFilter) ([]Tag, error) {
	path := fmt.Sprintf("/%s/tags", entityType)

	if filter != nil {
		path += "?"
		if filter.Limit > 0 {
			path += fmt.Sprintf("limit=%d&", filter.Limit)
		}
		if filter.Page > 0 {
			path += fmt.Sprintf("page=%d&", filter.Page)
		}
	}

	var resp TagsResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	return resp.Embedded.Tags, nil
}
