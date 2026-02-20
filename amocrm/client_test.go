package amocrm

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestNewClient(t *testing.T) {
	client := NewClient(
		WithSubdomain("test"),
		WithPermanentToken("token123"),
	)

	if client.subdomain != "test" {
		t.Errorf("Expected subdomain 'test', got '%s'", client.subdomain)
	}

	if client.permanentToken != "token123" {
		t.Errorf("Expected token 'token123', got '%s'", client.permanentToken)
	}

	if client.authType != AuthTypePermanentToken {
		t.Errorf("Expected auth type %d, got %d", AuthTypePermanentToken, client.authType)
	}
}

func TestClientWithOAuth2(t *testing.T) {
	client := NewClient(
		WithSubdomain("test"),
		WithOAuth2("client-id", "client-secret", "https://example.com/callback"),
	)

	if client.authType != AuthTypeOAuth2 {
		t.Errorf("Expected auth type %d, got %d", AuthTypeOAuth2, client.authType)
	}

	if client.oauth2Config == nil {
		t.Fatal("OAuth2 config should not be nil")
	}

	if client.oauth2Config.ClientID != "client-id" {
		t.Errorf("Expected client ID 'client-id', got '%s'", client.oauth2Config.ClientID)
	}
}

func TestAPIError(t *testing.T) {
	err := &APIError{
		StatusCode: 404,
		Message:    "Not Found",
	}

	expected := "API error (status 404): Not Found"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestTokenIsExpired(t *testing.T) {
	// Test expired token
	expiredToken := &Token{
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	if !expiredToken.IsExpired() {
		t.Error("Token should be expired")
	}

	// Test valid token
	validToken := &Token{
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if validToken.IsExpired() {
		t.Error("Token should not be expired")
	}
}

func TestContactsService_List(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v4/contacts" {
			t.Errorf("Expected path '/api/v4/contacts', got '%s'", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"_embedded": {
				"contacts": [
					{
						"id": 1,
						"name": "Test Contact"
					}
				]
			}
		}`))
	}))
	defer server.Close()

	// Create client with test server
	client := &Client{
		httpClient:     &http.Client{},
		subdomain:      "test",
		domain:         "amocrm.ru",
		baseURL:        server.URL + "/api/v4",
		authType:       AuthTypePermanentToken,
		permanentToken: "test-token",
		rateLimiter:    rate.NewLimiter(rate.Inf, 1),
		logger:         slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
	client.Contacts = &ContactsService{client: client}

	ctx := context.Background()
	contacts, err := client.Contacts.List(ctx, nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(contacts) != 1 {
		t.Errorf("Expected 1 contact, got %d", len(contacts))
	}

	if contacts[0].Name != "Test Contact" {
		t.Errorf("Expected contact name 'Test Contact', got '%s'", contacts[0].Name)
	}
}
