package amocrm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// AuthService handles OAuth 2.0 authentication
type AuthService struct {
	client *Client
}

// ExchangeCode exchanges an authorization code for access and refresh tokens
func (s *AuthService) ExchangeCode(ctx context.Context, code string) error {
	if s.client.authType != AuthTypeOAuth2 {
		return fmt.Errorf("OAuth2 is not configured")
	}

	if s.client.oauth2Config == nil {
		return fmt.Errorf("OAuth2 config is missing")
	}

	// Prepare request
	data := url.Values{}
	data.Set("client_id", s.client.oauth2Config.ClientID)
	data.Set("client_secret", s.client.oauth2Config.ClientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", s.client.oauth2Config.RedirectURI)

	tokenURL := fmt.Sprintf("https://%s.%s/oauth2/access_token", s.client.subdomain, s.client.domain)
	req, err := s.client.httpClient.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to exchange code: %w", err)
	}
	defer req.Body.Close()

	if req.StatusCode != 200 {
		return fmt.Errorf("token exchange failed with status %d", req.StatusCode)
	}

	var token Token
	if err := json.NewDecoder(req.Body).Decode(&token); err != nil {
		return fmt.Errorf("failed to decode token: %w", err)
	}

	token.ExpiresAt = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)

	// Save token
	s.client.tokenMu.Lock()
	s.client.currentToken = &token
	s.client.tokenMu.Unlock()

	// Persist token
	if s.client.tokenStorage != nil {
		domain := fmt.Sprintf("%s.%s", s.client.subdomain, s.client.domain)
		if err := s.client.tokenStorage.Save(ctx, domain, &token); err != nil {
			return fmt.Errorf("failed to save token: %w", err)
		}
	}

	return nil
}

// GetAuthorizationURL returns the OAuth2 authorization URL
func (s *AuthService) GetAuthorizationURL(state string, mode string) (string, error) {
	if s.client.authType != AuthTypeOAuth2 {
		return "", fmt.Errorf("OAuth2 is not configured")
	}

	if s.client.oauth2Config == nil {
		return "", fmt.Errorf("OAuth2 config is missing")
	}

	params := url.Values{}
	params.Set("client_id", s.client.oauth2Config.ClientID)
	params.Set("redirect_uri", s.client.oauth2Config.RedirectURI)
	params.Set("response_type", "code")

	if state != "" {
		params.Set("state", state)
	}

	if mode != "" {
		params.Set("mode", mode) // popup or post_message
	}

	authURL := fmt.Sprintf("https://%s.%s/oauth?%s", s.client.subdomain, s.client.domain, params.Encode())
	return authURL, nil
}

// RefreshToken manually refreshes the OAuth2 token
func (s *AuthService) RefreshToken(ctx context.Context) error {
	return s.client.refreshToken(ctx)
}

// GetCurrentToken returns the current OAuth2 token
func (s *AuthService) GetCurrentToken() *Token {
	s.client.tokenMu.RLock()
	defer s.client.tokenMu.RUnlock()
	return s.client.currentToken
}
