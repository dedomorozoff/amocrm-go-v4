package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ALipckin/amocrm-go-v4/amocrm"
)

// FileStorage implements TokenStorage using the file system
type FileStorage struct {
	directory string
}

// NewFileStorage creates a new file-based token storage
func NewFileStorage(directory string) *FileStorage {
	return &FileStorage{
		directory: directory,
	}
}

// Save saves a token to a file
func (s *FileStorage) Save(ctx context.Context, domain string, token *amocrm.Token) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(s.directory, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal token to JSON
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	// Write to file
	filename := filepath.Join(s.directory, domain+".json")
	if err := os.WriteFile(filename, data, 0600); err != nil {
		return fmt.Errorf("failed to write token file: %w", err)
	}

	return nil
}

// Load loads a token from a file
func (s *FileStorage) Load(ctx context.Context, domain string) (*amocrm.Token, error) {
	filename := filepath.Join(s.directory, domain+".json")

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, nil
	}

	// Read file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read token file: %w", err)
	}

	// Unmarshal token
	var token amocrm.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	return &token, nil
}

// HasToken checks if a token file exists
func (s *FileStorage) HasToken(ctx context.Context, domain string) (bool, error) {
	filename := filepath.Join(s.directory, domain+".json")
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
