package amocrm

import (
	"context"
	"fmt"
)

// TaskType represents task type constants
type TaskType int

const (
	TaskTypeCall TaskType = 1
	TaskTypeMeet TaskType = 2
	TaskTypeMail TaskType = 3
)

// Task represents an AmoCRM task
type Task struct {
	ID                int         `json:"id,omitempty"`
	CreatedBy         int         `json:"created_by,omitempty"`
	UpdatedBy         int         `json:"updated_by,omitempty"`
	CreatedAt         int64       `json:"created_at,omitempty"`
	UpdatedAt         int64       `json:"updated_at,omitempty"`
	ResponsibleUserID int         `json:"responsible_user_id,omitempty"`
	GroupID           int         `json:"group_id,omitempty"`
	EntityID          int         `json:"entity_id,omitempty"`
	EntityType        string      `json:"entity_type,omitempty"` // leads, contacts, companies, customers
	IsCompleted       bool        `json:"is_completed,omitempty"`
	TaskTypeID        int         `json:"task_type_id,omitempty"`
	Text              string      `json:"text"`
	Duration          int         `json:"duration,omitempty"`
	CompleteTill      int64       `json:"complete_till"`
	Result            *TaskResult `json:"result,omitempty"`
	AccountID         int         `json:"account_id,omitempty"`
}

// TaskResult represents task completion result
type TaskResult struct {
	Text string `json:"text,omitempty"`
}

// TasksService handles communication with task-related methods
type TasksService struct {
	client *Client
}

// TasksResponse represents the API response for tasks list
type TasksResponse struct {
	Embedded struct {
		Tasks []Task `json:"tasks"`
	} `json:"_embedded"`
	Links Links `json:"_links"`
	Page  Page  `json:"_page,omitempty"`
}

// TasksFilter represents filter options for listing tasks
type TasksFilter struct {
	Limit             int
	Page              int
	Filter            map[string]interface{}
	Order             string
	ResponsibleUserID int
	IsCompleted       *bool
}

// List retrieves a list of tasks
func (s *TasksService) List(ctx context.Context, filter *TasksFilter) ([]Task, error) {
	path := "/tasks"

	if filter != nil {
		path += "?"
		if filter.Limit > 0 {
			path += fmt.Sprintf("limit=%d&", filter.Limit)
		}
		if filter.Page > 0 {
			path += fmt.Sprintf("page=%d&", filter.Page)
		}
		if filter.ResponsibleUserID > 0 {
			path += fmt.Sprintf("filter[responsible_user_id]=%d&", filter.ResponsibleUserID)
		}
		if filter.IsCompleted != nil {
			completed := 0
			if *filter.IsCompleted {
				completed = 1
			}
			path += fmt.Sprintf("filter[is_completed]=%d&", completed)
		}
	}

	var resp TasksResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	return resp.Embedded.Tasks, nil
}

// GetByID retrieves a task by ID
func (s *TasksService) GetByID(ctx context.Context, id int) (*Task, error) {
	path := fmt.Sprintf("/tasks/%d", id)

	var task Task
	if err := s.client.GetJSON(ctx, path, &task); err != nil {
		return nil, err
	}

	return &task, nil
}

// Create creates a new task
func (s *TasksService) Create(ctx context.Context, task *Task) (*Task, error) {
	type request struct {
		Tasks []Task `json:"tasks"`
	}

	req := request{
		Tasks: []Task{*task},
	}

	var resp TasksResponse
	if err := s.client.PostJSON(ctx, "/tasks", req, &resp); err != nil {
		return nil, err
	}

	if len(resp.Embedded.Tasks) == 0 {
		return nil, fmt.Errorf("no task returned from API")
	}

	return &resp.Embedded.Tasks[0], nil
}

// CreateBatch creates multiple tasks in one request
func (s *TasksService) CreateBatch(ctx context.Context, tasks []*Task) ([]Task, error) {
	type request struct {
		Tasks []Task `json:"tasks"`
	}

	tasksValues := make([]Task, len(tasks))
	for i, t := range tasks {
		tasksValues[i] = *t
	}

	req := request{
		Tasks: tasksValues,
	}

	var resp TasksResponse
	if err := s.client.PostJSON(ctx, "/tasks", req, &resp); err != nil {
		return nil, err
	}

	return resp.Embedded.Tasks, nil
}

// Update updates an existing task
func (s *TasksService) Update(ctx context.Context, task *Task) (*Task, error) {
	if task.ID == 0 {
		return nil, fmt.Errorf("task ID is required for update")
	}

	type request struct {
		Tasks []Task `json:"tasks"`
	}

	req := request{
		Tasks: []Task{*task},
	}

	var resp TasksResponse
	if err := s.client.PatchJSON(ctx, "/tasks", req, &resp); err != nil {
		return nil, err
	}

	if len(resp.Embedded.Tasks) == 0 {
		return nil, fmt.Errorf("no task returned from API")
	}

	return &resp.Embedded.Tasks[0], nil
}

// Complete marks a task as completed
func (s *TasksService) Complete(ctx context.Context, taskID int, resultText string) error {
	task := &Task{
		ID:          taskID,
		IsCompleted: true,
		Result: &TaskResult{
			Text: resultText,
		},
	}

	_, err := s.Update(ctx, task)
	return err
}
