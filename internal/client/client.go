// Package client provides the client for making HTTP requests.
package client

//go:generate mockgen -destination=mocks/client.go -package=mock_client -source=client.go

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/twk/skeleton-go-api/internal/logger"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// AuthType represents the type of authentication to use.
type AuthType int

const (
	// AuthTypeToken represents the token authentication type.
	AuthTypeToken = iota
	// AuthTypeBearer represents the bearer authentication type.
	AuthTypeBearer
	// AuthTypeBasic represents the basic authentication type.
	AuthTypeBasic
)

// Client is a wrapper around the http client.
type Client struct {
	httpClient httpClient
}

// NewClient creates a new Client.
func NewClient(httpClient httpClient) (*Client, error) {
	return &Client{
		httpClient: httpClient,
	}, nil
}

// Request performs an HTTP request with the specified method, URL, headers, query parameters, and body.
func (c *Client) Request(ctx context.Context, _ *logger.Logger, method, targetURL, path string, headers, query map[string]string, body io.Reader) (*http.Response, error) {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	refURL, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse path: %w", err)
	}

	urlPath := parsedURL.ResolveReference(refURL).String()

	req, err := http.NewRequestWithContext(ctx, method, urlPath, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.SetHeaders(req, headers)

	q := req.URL.Query()
	for key, val := range query {
		q.Add(key, val)
	}

	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}

	return resp, nil
}

// SetHeaders sets the headers for the request
func (c *Client) SetHeaders(req *http.Request, headers map[string]string) {
	for key, val := range headers {
		req.Header.Set(key, val)
	}
}
