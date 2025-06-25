package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"ticktick-tui/internal/models"
)

const (
	AuthURL  = "https://ticktick.com/oauth/authorize"
	TokenURL = "https://ticktick.com/oauth/token"
)

type OAuthClient struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// GetAuthURL generates the authorization URL
func (c *OAuthClient) GetAuthURL() string {
	params := url.Values{}
	params.Set("client_id", c.ClientID)
	params.Set("scope", "tasks:read tasks:write")
	params.Set("state", "ticktick-tui-state")
	params.Set("redirect_uri", c.RedirectURI)
	params.Set("response_type", "code")

	return AuthURL + "?" + params.Encode()
}

// ExchangeCodeForToken exchanges authorization code for access token
func (c *OAuthClient) ExchangeCodeForToken(code, scope string) (*models.OAuthToken, error) {
	data := url.Values{}
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("scope", scope)
	data.Set("redirect_uri", c.RedirectURI)

	req, err := http.NewRequest("POST", TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Basic Auth
	auth := base64.StdEncoding.EncodeToString([]byte(c.ClientID + ":" + c.ClientSecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token request failed: %s", string(body))
	}

	var token models.OAuthToken
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, fmt.Errorf("parsing token response: %w", err)
	}

	return &token, nil
}
