package amocrm

import (
	"context"
)

type Pipeline struct {
	ID           int              `json:"id"`
	Name         string           `json:"name"`
	Sort         int              `json:"sort"`
	IsMain       bool             `json:"is_main"`
	IsUnsortedOn bool             `json:"is_unsorted_on,omitempty"`
	IsArchive    bool             `json:"is_archive,omitempty"`
	AccountID    int              `json:"account_id,omitempty"`
	Embedded     PipelineEmbedded `json:"_embedded,omitempty"`
}

type PipelineEmbedded struct {
	Statuses []Status `json:"statuses,omitempty"`
}

type Status struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Sort       int    `json:"sort"`
	IsEditable bool   `json:"is_editable"`
	PipelineID int    `json:"pipeline_id"`
	Color      string `json:"color,omitempty"`
	Type       int    `json:"type,omitempty"`
	AccountID  int    `json:"account_id,omitempty"`
}

type PipelinesService struct {
	client *Client
}

type PipelinesResponse struct {
	Embedded struct {
		Pipelines []Pipeline `json:"pipelines"`
	} `json:"_embedded"`
	Links Links `json:"_links"`
}

func (s *PipelinesService) List(ctx context.Context) ([]Pipeline, error) {
	path := "/leads/pipelines"

	var resp PipelinesResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	return resp.Embedded.Pipelines, nil
}
