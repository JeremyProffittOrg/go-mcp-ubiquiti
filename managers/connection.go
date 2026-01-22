// Package managers provides business logic for UniFi API interactions.
package managers

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"

	"github.com/sirkirby/go-mcp-ubiquiti/config"
)

// ConnectionManager handles authentication and API requests to the UniFi controller.
type ConnectionManager struct {
	config    *config.Config
	client    *http.Client
	baseURL   string
	isUnifiOS bool
	apiPrefix string
	csrfToken string
	mu        sync.RWMutex
	loggedIn  bool
}

// APIResponse represents a standard UniFi API response.
type APIResponse struct {
	Meta struct {
		RC  string `json:"rc"`
		Msg string `json:"msg,omitempty"`
	} `json:"meta"`
	Data json.RawMessage `json:"data"`
}

// NewConnectionManager creates a new connection manager.
func NewConnectionManager(cfg *config.Config) *ConnectionManager {
	jar, _ := cookiejar.New(nil)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: !cfg.VerifySSL,
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
		Jar:       jar,
	}

	return &ConnectionManager{
		config:  cfg,
		client:  client,
		baseURL: cfg.BaseURL(),
	}
}

// Login authenticates with the UniFi controller.
func (c *ConnectionManager) Login(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Detect controller type if auto
	if c.config.ControllerType == "auto" {
		if err := c.detectControllerType(ctx); err != nil {
			log.Printf("Controller type detection failed, assuming UniFi OS: %v", err)
			c.isUnifiOS = true
		}
	} else {
		c.isUnifiOS = c.config.ControllerType == "proxy"
	}

	// Set API prefix based on controller type
	if c.isUnifiOS {
		c.apiPrefix = "/proxy/network"
	} else {
		c.apiPrefix = ""
	}

	// Perform login
	loginPath := "/api/login"
	if c.isUnifiOS {
		loginPath = "/api/auth/login"
	}

	credentials := map[string]string{
		"username": c.config.Username,
		"password": c.config.Password,
	}

	body, err := json.Marshal(credentials)
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+loginPath, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read login response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// Extract CSRF token from cookies if present
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "csrf_token" || cookie.Name == "X-CSRF-Token" {
			c.csrfToken = cookie.Value
		}
	}

	// Also check response header for CSRF token
	if token := resp.Header.Get("X-CSRF-Token"); token != "" {
		c.csrfToken = token
	}

	c.loggedIn = true
	log.Printf("Successfully logged in to UniFi controller (UniFi OS: %v)", c.isUnifiOS)

	return nil
}

// detectControllerType probes the controller to determine its type.
func (c *ConnectionManager) detectControllerType(ctx context.Context) error {
	// Try UniFi OS login endpoint first
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/auth/login", nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	// UniFi OS typically returns 200 or 401 on the auth endpoint
	// Standalone returns 404
	if resp.StatusCode != http.StatusNotFound {
		c.isUnifiOS = true
		return nil
	}

	// Try standalone login endpoint
	req, err = http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/login", nil)
	if err != nil {
		return err
	}

	resp, err = c.client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		c.isUnifiOS = false
		return nil
	}

	return fmt.Errorf("unable to detect controller type")
}

// Request makes an authenticated API request.
func (c *ConnectionManager) Request(ctx context.Context, method, path string, body interface{}) (json.RawMessage, error) {
	c.mu.RLock()
	if !c.loggedIn {
		c.mu.RUnlock()
		if err := c.Login(ctx); err != nil {
			return nil, fmt.Errorf("not logged in and login failed: %w", err)
		}
		c.mu.RLock()
	}
	c.mu.RUnlock()

	data, err := c.doRequest(ctx, method, path, body)
	if err != nil {
		// Check if it's an auth error and retry
		if isAuthError(err) {
			log.Println("Session expired, re-authenticating...")
			if loginErr := c.Login(ctx); loginErr != nil {
				return nil, fmt.Errorf("re-authentication failed: %w", loginErr)
			}
			return c.doRequest(ctx, method, path, body)
		}
		return nil, err
	}

	return data, nil
}

func (c *ConnectionManager) doRequest(ctx context.Context, method, path string, body interface{}) (json.RawMessage, error) {
	fullPath := c.buildPath(path)

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+fullPath, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.csrfToken != "" {
		req.Header.Set("X-CSRF-Token", c.csrfToken)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, &AuthError{Message: "unauthorized"}
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse the response
	var apiResp APIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		// Some endpoints return raw data without the meta wrapper
		return respBody, nil
	}

	if apiResp.Meta.RC == "error" {
		if apiResp.Meta.Msg == "api.err.LoginRequired" {
			return nil, &AuthError{Message: apiResp.Meta.Msg}
		}
		return nil, fmt.Errorf("API error: %s", apiResp.Meta.Msg)
	}

	return apiResp.Data, nil
}

// buildPath constructs the full API path based on controller type.
func (c *ConnectionManager) buildPath(endpoint string) string {
	return c.apiPrefix + endpoint
}

// GetSitePath returns the site-specific API path.
func (c *ConnectionManager) GetSitePath(endpoint string) string {
	return fmt.Sprintf("/api/s/%s%s", c.config.Site, endpoint)
}

// IsUnifiOS returns whether the controller is UniFi OS based.
func (c *ConnectionManager) IsUnifiOS() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.isUnifiOS
}

// Site returns the configured site name.
func (c *ConnectionManager) Site() string {
	return c.config.Site
}

// AuthError represents an authentication error.
type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return "authentication error: " + e.Message
}

func isAuthError(err error) bool {
	_, ok := err.(*AuthError)
	return ok
}
