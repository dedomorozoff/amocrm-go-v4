// Package amocrm provides a Go client for the AmoCRM API v4.
//
// This library supports:
//   - OAuth 2.0 authorization with automatic token refresh
//   - Permanent tokens (recommended for server integrations)
//   - Rate limiting (7 requests per second by default)
//   - Context support for timeouts and cancellation
//   - Type-safe API interactions
//
// Example usage with permanent token:
//
//	client := amocrm.NewClient(
//		amocrm.WithSubdomain("testsubdomain"),
//		amocrm.WithPermanentToken("your-permanent-token"),
//	)
//
//	ctx := context.Background()
//	account, err := client.Account.Get(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Example usage with OAuth 2.0:
//
//	storage := storage.NewFileStorage("./tokens")
//	client := amocrm.NewClient(
//		amocrm.WithSubdomain("testsubdomain"),
//		amocrm.WithOAuth2("client-id", "client-secret", "redirect-uri"),
//		amocrm.WithTokenStorage(storage),
//	)
//
//	err := client.Auth.ExchangeCode(ctx, "auth-code")
package amocrm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	// DefaultDomain is the default AmoCRM domain
	DefaultDomain = "amocrm.ru"

	// DefaultRateLimit is the default rate limit (requests per second)
	DefaultRateLimit = 7

	// DefaultTimeout is the default HTTP client timeout
	DefaultTimeout = 30 * time.Second

	// APIVersion is the AmoCRM API version
	APIVersion = "v4"
)

// Client is the main AmoCRM API client
type Client struct {
	// HTTP client
	httpClient *http.Client

	// Configuration
	subdomain string
	domain    string
	baseURL   string

	// Authentication
	authType       AuthType
	permanentToken string
	oauth2Config   *OAuth2Config
	tokenStorage   TokenStorage
	currentToken   *Token
	tokenMu        sync.RWMutex

	// Rate limiting
	rateLimiter *rate.Limiter

	// Logging
	logger *slog.Logger
	debug  bool

	// API Services
	Account    *AccountService
	Contacts   *ContactsService
	Companies  *CompaniesService
	Leads      *LeadsService
	Tasks      *TasksService
	Notes      *NotesService
	Webhooks   *WebhooksService
	Catalogs   *CatalogsService
	Auth       *AuthService
	Users      *UserService
	Roles      *RoleService
	Pipelines  *PipelinesService
	TaskTypes  *TaskTypesService
	Tags       *TagsService
	Events     *EventsService
	Pagination *PaginationService
}

// AuthType represents the type of authentication
type AuthType int

const (
	AuthTypeNone AuthType = iota
	AuthTypePermanentToken
	AuthTypeOAuth2
)

// OAuth2Config holds OAuth 2.0 configuration
type OAuth2Config struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// Token represents an OAuth 2.0 token
type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// IsExpired checks if the token is expired
func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// ClientOption is a function that configures the Client
type ClientOption func(*Client)

// WithSubdomain sets the AmoCRM subdomain
func WithSubdomain(subdomain string) ClientOption {
	return func(c *Client) {
		c.subdomain = subdomain
	}
}

// WithDomain sets the AmoCRM domain (default: amocrm.ru)
func WithDomain(domain string) ClientOption {
	return func(c *Client) {
		c.domain = domain
	}
}

// WithPermanentToken sets permanent token authentication
func WithPermanentToken(token string) ClientOption {
	return func(c *Client) {
		c.authType = AuthTypePermanentToken
		c.permanentToken = token
	}
}

// WithOAuth2 sets OAuth 2.0 authentication
func WithOAuth2(clientID, clientSecret, redirectURI string) ClientOption {
	return func(c *Client) {
		c.authType = AuthTypeOAuth2
		c.oauth2Config = &OAuth2Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURI:  redirectURI,
		}
	}
}

// WithTokenStorage sets the token storage implementation
func WithTokenStorage(storage TokenStorage) ClientOption {
	return func(c *Client) {
		c.tokenStorage = storage
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithRateLimit sets the rate limit (requests per second)
func WithRateLimit(rps int) ClientOption {
	return func(c *Client) {
		c.rateLimiter = rate.NewLimiter(rate.Limit(rps), 1)
	}
}

// WithTimeout sets the HTTP client timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithLogger sets a custom logger
func WithLogger(logger *slog.Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// WithDebug enables debug logging
func WithDebug(debug bool) ClientOption {
	return func(c *Client) {
		c.debug = debug
	}
}

// NewClient creates a new AmoCRM API client
func NewClient(opts ...ClientOption) *Client {
	client := &Client{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		domain:      DefaultDomain,
		rateLimiter: rate.NewLimiter(rate.Limit(DefaultRateLimit), 1),
		logger:      slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	// Apply options
	for _, opt := range opts {
		opt(client)
	}

	// Validate configuration
	if client.subdomain == "" {
		panic("subdomain is required")
	}

	// Build base URL
	client.baseURL = fmt.Sprintf("https://%s.%s/api/%s", client.subdomain, client.domain, APIVersion)

	// Initialize services
	client.Account = &AccountService{client: client}
	client.Contacts = &ContactsService{client: client}
	client.Companies = &CompaniesService{client: client}
	client.Leads = &LeadsService{client: client}
	client.Tasks = &TasksService{client: client}
	client.Notes = &NotesService{client: client}
	client.Webhooks = &WebhooksService{client: client}
	client.Catalogs = &CatalogsService{client: client}
	client.Auth = &AuthService{client: client}
	client.Users = &UserService{client: client}
	client.Roles = &RoleService{client: client}
	client.Pipelines = &PipelinesService{client: client}
	client.TaskTypes = &TaskTypesService{client: client}
	client.Tags = &TagsService{client: client}
	client.Events = &EventsService{client: client}
	client.Pagination = &PaginationService{client: client}

	// Load token if using OAuth2
	if client.authType == AuthTypeOAuth2 && client.tokenStorage != nil {
		domain := fmt.Sprintf("%s.%s", client.subdomain, client.domain)
		token, err := client.tokenStorage.Load(context.Background(), domain)
		if err == nil && token != nil {
			client.currentToken = token
		}
	}

	return client
}

// do executes an HTTP request with rate limiting and authentication
func (c *Client) do(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	// Wait for rate limiter
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter error: %w", err)
	}

	// Build URL
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "amocrm-go/1.0")

	// Add authentication
	if err := c.addAuth(ctx, req); err != nil {
		return nil, fmt.Errorf("authentication error: %w", err)
	}

	// Log request if debug is enabled
	if c.debug {
		c.logger.Debug("API Request",
			"method", method,
			"url", u.String(),
		)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Log response if debug is enabled
	if c.debug {
		c.logger.Debug("API Response",
			"status", resp.StatusCode,
			"url", u.String(),
		)
	}

	// Handle 401 Unauthorized - try to refresh token
	if resp.StatusCode == http.StatusUnauthorized && c.authType == AuthTypeOAuth2 {
		resp.Body.Close()
		if err := c.refreshToken(ctx); err != nil {
			return nil, fmt.Errorf("token refresh failed: %w", err)
		}
		// Retry request with new token
		return c.do(ctx, method, path, body)
	}

	// Check for API errors
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(bodyBytes),
		}
	}

	return resp, nil
}

// addAuth adds authentication to the request
func (c *Client) addAuth(ctx context.Context, req *http.Request) error {
	switch c.authType {
	case AuthTypePermanentToken:
		req.Header.Set("Authorization", "Bearer "+c.permanentToken)
		return nil

	case AuthTypeOAuth2:
		c.tokenMu.RLock()
		token := c.currentToken
		c.tokenMu.RUnlock()

		if token == nil {
			return fmt.Errorf("no OAuth2 token available")
		}

		// Check if token is expired
		if token.IsExpired() {
			c.tokenMu.RUnlock()
			if err := c.refreshToken(ctx); err != nil {
				return err
			}
			c.tokenMu.RLock()
			token = c.currentToken
			c.tokenMu.RUnlock()
		}

		req.Header.Set("Authorization", "Bearer "+token.AccessToken)
		return nil

	default:
		return fmt.Errorf("no authentication method configured")
	}
}

// refreshToken refreshes the OAuth2 token
func (c *Client) refreshToken(ctx context.Context) error {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()

	if c.currentToken == nil || c.currentToken.RefreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	// Prepare request
	data := url.Values{}
	data.Set("client_id", c.oauth2Config.ClientID)
	data.Set("client_secret", c.oauth2Config.ClientSecret)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", c.currentToken.RefreshToken)
	data.Set("redirect_uri", c.oauth2Config.RedirectURI)

	tokenURL := fmt.Sprintf("https://%s.%s/oauth2/access_token", c.subdomain, c.domain)
	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("token refresh failed: %s", string(bodyBytes))
	}

	var token Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return err
	}

	token.ExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	c.currentToken = &token

	// Save token
	if c.tokenStorage != nil {
		domain := fmt.Sprintf("%s.%s", c.subdomain, c.domain)
		if err := c.tokenStorage.Save(ctx, domain, &token); err != nil {
			c.logger.Warn("Failed to save token", "error", err)
		}
	}

	return nil
}

// GetJSON performs a GET request and decodes JSON response
func (c *Client) GetJSON(ctx context.Context, path string, result interface{}) error {
	resp, err := c.do(ctx, "GET", path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

// PostJSON performs a POST request with JSON body
func (c *Client) PostJSON(ctx context.Context, path string, body interface{}, result interface{}) error {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := c.do(ctx, "POST", path, strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}

	return nil
}

// PatchJSON performs a PATCH request with JSON body
func (c *Client) PatchJSON(ctx context.Context, path string, body interface{}, result interface{}) error {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := c.do(ctx, "PATCH", path, strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if result != nil {
		return json.NewDecoder(resp.Body).Decode(result)
	}

	return nil
}

// DeleteJSON performs a DELETE request
func (c *Client) DeleteJSON(ctx context.Context, path string) error {
	resp, err := c.do(ctx, "DELETE", path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
