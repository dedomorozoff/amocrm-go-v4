package amocrm

import (
	"context"
)

type TaskTypeItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TaskTypesService struct {
	client *Client
}

type TaskTypesResponse struct {
	Embedded struct {
		TaskTypes []TaskTypeItem `json:"task_types"`
	} `json:"_embedded"`
	Links Links `json:"_links"`
}

func (s *TaskTypesService) List(ctx context.Context) ([]TaskTypeItem, error) {
	// Task types are predefined in AmoCRM
	taskTypes := []TaskTypeItem{
		{ID: 1, Name: "Звонок"},
		{ID: 2, Name: "Встреча"},
		{ID: 3, Name: "Написать письмо"},
	}
	return taskTypes, nil
}
