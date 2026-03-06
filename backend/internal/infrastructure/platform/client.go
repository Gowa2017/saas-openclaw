package platform

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Common errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrPlatformUnavailable = errors.New("platform service unavailable")
	ErrRequestTimeout    = errors.New("request timeout")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrInvalidResponse   = errors.New("invalid response from platform")
)

// UserInfo represents the user information from the platform
type UserInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	TenantID string `json:"tenantId"`
	Role     string `json:"role"`
}

// Client is the platform API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
}

// ClientOption is a functional option for configuring the client
type ClientOption func(*Client)

// WithTimeout sets the request timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.timeout = timeout
		c.httpClient.Timeout = timeout
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// NewClient creates a new platform client
func NewClient(baseURL string, opts ...ClientOption) *Client {
	c := &Client{
		baseURL: baseURL,
		timeout: 5 * time.Second,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// GetUserInfo fetches user information from the platform
func (c *Client) GetUserInfo(ctx context.Context, token string) (*UserInfo, error) {
	url := fmt.Sprintf("%s/api/user/info", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, ErrRequestTimeout
		}
		return nil, fmt.Errorf("%w: %v", ErrPlatformUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrUserNotFound
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%w: status %d, body: %s", ErrPlatformUnavailable, resp.StatusCode, string(body))
	}

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidResponse, err)
	}

	if userInfo.ID == "" {
		return nil, fmt.Errorf("%w: missing user id", ErrInvalidResponse)
	}

	return &userInfo, nil
}

// HealthCheck checks if the platform service is healthy
func (c *Client) HealthCheck(ctx context.Context) error {
	url := fmt.Sprintf("%s/health", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrPlatformUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: status %d", ErrPlatformUnavailable, resp.StatusCode)
	}

	return nil
}
