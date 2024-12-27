// Package client provides the client for making HTTP requests.
package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// AuthType represents the type of authentication to use.
type AuthType int

const (
	AuthTypeToken = iota
	AuthTypeBearer
	AuthTypeBasic
)

// Client is a wrapper around the http client.
type Client struct {
	authType   AuthType
	baseURL    *url.URL
	httpClient httpClient
	token      string
}

// NewClient creates a new Client.
func NewClient(baseURL, token string, authType AuthType, httpClient httpClient) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	return &Client{
		baseURL:    parsedURL,
		authType:   authType,
		token:      token,
		httpClient: httpClient,
	}, nil
}

// Get performs a GET request to the specified URL with the specified query parameters. I need to encode space as %20 for query parameters, not +
func (c *Client) Get(ctx context.Context, urlPath string, query map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlPath, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.SetHeaders(req, map[string]string{
		"Accept": "application/json",
	})

	c.SetAuthType(req)

	q := req.URL.Query()
	for key, val := range query {
		q.Add(key, val)
	}

	req.URL.RawQuery = strings.ReplaceAll(q.Encode(), "+", "%20")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}

	return resp, nil
}

// SetAuthType sets the authentication type for the client
func (c *Client) SetAuthType(req *http.Request) {
	switch c.authType {
	case AuthTypeToken:
		req.Header.Set("Authorization", fmt.Sprintf("token %s", c.token))
	case AuthTypeBearer:
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	case AuthTypeBasic:
		req.SetBasicAuth(c.token, "")
	}
}

// SetHeaders sets the headers for the request
func (c *Client) SetHeaders(req *http.Request, headers map[string]string) {
	for key, val := range headers {
		req.Header.Set(key, val)
	}
}
