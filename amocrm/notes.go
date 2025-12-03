package amocrm

import (
	"context"
	"fmt"
)

// NoteType represents note type constants
type NoteType string

const (
	NoteTypeCommon         NoteType = "common"
	NoteTypeCallIn         NoteType = "call_in"
	NoteTypeCallOut        NoteType = "call_out"
	NoteTypeSMSIn          NoteType = "sms_in"
	NoteTypeSMSOut         NoteType = "sms_out"
	NoteTypeServiceMessage NoteType = "service_message"
)

// Note represents an AmoCRM note
type Note struct {
	ID                int                    `json:"id,omitempty"`
	EntityID          int                    `json:"entity_id"`
	CreatedBy         int                    `json:"created_by,omitempty"`
	UpdatedBy         int                    `json:"updated_by,omitempty"`
	CreatedAt         int64                  `json:"created_at,omitempty"`
	UpdatedAt         int64                  `json:"updated_at,omitempty"`
	ResponsibleUserID int                    `json:"responsible_user_id,omitempty"`
	GroupID           int                    `json:"group_id,omitempty"`
	NoteType          NoteType               `json:"note_type"`
	Params            map[string]interface{} `json:"params"`
	AccountID         int                    `json:"account_id,omitempty"`
}

// NotesService handles communication with note-related methods
type NotesService struct {
	client *Client
}

// NotesResponse represents the API response for notes list
type NotesResponse struct {
	Embedded struct {
		Notes []Note `json:"notes"`
	} `json:"_embedded"`
	Links Links `json:"_links"`
	Page  Page  `json:"_page,omitempty"`
}

// NotesFilter represents filter options for listing notes
type NotesFilter struct {
	Limit      int
	Page       int
	NoteType   []NoteType
	EntityID   int
	EntityType EntityType
}

// List retrieves a list of notes for an entity
func (s *NotesService) List(ctx context.Context, entityType EntityType, entityID int, filter *NotesFilter) ([]Note, error) {
	path := fmt.Sprintf("/%s/%d/notes", entityType, entityID)

	if filter != nil {
		path += "?"
		if filter.Limit > 0 {
			path += fmt.Sprintf("limit=%d&", filter.Limit)
		}
		if filter.Page > 0 {
			path += fmt.Sprintf("page=%d&", filter.Page)
		}
		for _, noteType := range filter.NoteType {
			path += fmt.Sprintf("filter[note_type][]=%s&", noteType)
		}
	}

	var resp NotesResponse
	if err := s.client.GetJSON(ctx, path, &resp); err != nil {
		return nil, err
	}

	return resp.Embedded.Notes, nil
}

// GetByID retrieves a note by ID
func (s *NotesService) GetByID(ctx context.Context, entityType EntityType, entityID int, noteID int) (*Note, error) {
	path := fmt.Sprintf("/%s/%d/notes/%d", entityType, entityID, noteID)

	var note Note
	if err := s.client.GetJSON(ctx, path, &note); err != nil {
		return nil, err
	}

	return &note, nil
}

// Create creates a new note
func (s *NotesService) Create(ctx context.Context, entityType EntityType, note *Note) (*Note, error) {
	type request struct {
		Notes []Note `json:"notes"`
	}

	req := request{
		Notes: []Note{*note},
	}

	path := fmt.Sprintf("/%s/%d/notes", entityType, note.EntityID)

	var resp NotesResponse
	if err := s.client.PostJSON(ctx, path, req, &resp); err != nil {
		return nil, err
	}

	if len(resp.Embedded.Notes) == 0 {
		return nil, fmt.Errorf("no note returned from API")
	}

	return &resp.Embedded.Notes[0], nil
}

// CreateBatch creates multiple notes in one request
func (s *NotesService) CreateBatch(ctx context.Context, entityType EntityType, entityID int, notes []*Note) ([]Note, error) {
	type request struct {
		Notes []Note `json:"notes"`
	}

	notesValues := make([]Note, len(notes))
	for i, n := range notes {
		notesValues[i] = *n
	}

	req := request{
		Notes: notesValues,
	}

	path := fmt.Sprintf("/%s/%d/notes", entityType, entityID)

	var resp NotesResponse
	if err := s.client.PostJSON(ctx, path, req, &resp); err != nil {
		return nil, err
	}

	return resp.Embedded.Notes, nil
}
