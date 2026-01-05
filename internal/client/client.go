// Package client provides the client for making HTTP requests.
package client

import (
	"context"
	"fmt"
	"net/http"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type AuthType int

const (
	NoAuth AuthType = iota
	BasicAuth
	BearerToken
	Token
)

// Client is a wrapper around the http client.
type Client struct {
	authType  AuthType
	baseURL *url.URL
	httpClient httpClient
	token    string
	logger  Logger
}

type paginatedResponse[T any] struct {
	Items      []T   `json:"items"`
	NextPage   string `json:"next_page"`
	TotalCount int    `json:"total_count"`
}

// NewClient creates a new Client.
func NewClient(baseURL, token string, authType AuthType, httpClient httpClient, logger Logger) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	return &Client{
		authType:  authType,
		baseURL:   parsedURL,
		httpClient: httpClient,
		token:     token,
		logger:    logger,
	}
}

// Get makes a GET request to the specified path.
func (c *Client) Get(ctx context.Context, path string, queryParams map[string]string) (*http.Response, error) {
	targetURL, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse path: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.SetHeaders(req, map[string]string{
		"Accept": "application/json",
	})
	c.SetAuth(req)

	// Add query parameters
	q := req.URL.Query()
	for key, value := range queryParams {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// Post makes a POST request to the specified path.
func (c *Client) Post(ctx context.Context, path string, body io.Reader, headers map[string]string) (*http.Response, error) {
	targetURL, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse path: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.SetHeaders(req, headers)
	c.SetAuth(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// SetAuth sets the appropriate authentication headers.
func (c *Client) SetAuth(req *http.Request) {
	switch c.authType {
	case BasicAuth:
		req.SetBasicAuth("user", c.token)
	case BearerToken:
		req.Header.Set("Authorization", "Bearer "+c.token)
	case Token:
		req.Header.Set("X-Auth-Token", c.token)
	}
}

// SetHeader sets a header for the request.
func (c *Client) SetHeader(req *http.Request, header map[string]string) {
	for key, value := range header {
		req.Header.Set(key, value)
	}
}

func get[T any](ctx context.Context, client *Client, path string, queryParams map[string]string) (*T, error) {
	resp, err := client.Get(ctx, path, queryParams)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(respBody))
	}

	var result T
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &result, nil
}

func Post[T any, B any](ctx context.Context, client *Client, path string, body *B, headers map[string]string) (*T, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	
	bodyReader := bytes.NewReader(bodyBytes)

	resp, err := client.Post(ctx, path, bodyReader, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(respBody))
	}

	var result T
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &result, nil
}

func pagedGet[T any](ctx context.Context, client *Client, path string, queryParams map[string]string) (*paginatedResponse[T], error) {
	resp, err := client.Get(ctx, path, queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to make paged GET request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(respBody))
	}

	var pagedResp paginatedResponse[T]
	if err := json.Unmarshal(respBody, &pagedResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal paged response body: %w", err)
	}

	return &pagedResp, nil
}

func paginate[T any](ctx context.Context, client *Client, path string, query map[string]string, getPageFunc func(context.Context, *Client, string, map[string]string) (*paginatedResponse[T], error)) ([]T, error) {
	results := make([]T, 0)
	nextPage := ""

	for {
		if nextPage != "" {
			query["page"] = nextPage
		}

		pageResp, err := getPageFunc(ctx, client, path, query)
		if err != nil {
			return nil, fmt.Errorf("failed to get page: %w", err)
		}

		results = append(results, pageResp.Items...)

		if pageResp.NextPage == "" {
			break
		}

		nextPageURL, err := url.Parse(pageResp.NextPage)
		if err != nil {
			return nil, fmt.Errorf("failed to parse next page URL: %w", err)
		}
		nextPage = nextPageURL.Query().Get("page")
	}

	return results, nil
}

// Example of how to use paginate
// func ExamplePaginateUsers(ctx context.Context, client *Client) ([]User, error) {
// 	return paginate[User](ctx, client, "/users", map[string]string{}, pagedGet[User])
// }

func pagedGetGeneric[R any](ctx context.Context, client *Client, path string, queryParams map[string]string) (*R, error) {
	resp, err := client.Get(ctx, path, queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to make paged GET request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, string(respBody))
	}

	var pagedResp R
	if err := json.Unmarshal(respBody, &pagedResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal paged response body: %w", err)
	}

	return &pagedResp, nil
}

func paginateGeneric[T any, R any](
	ctx context.Context,
	client *Client,
	path string,
	query map[string]string,
	getPageFunc func(context.Context, *Client, string, map[string]string) (*R, error),
	getItems func(*R) []T,
	getNextPage func(*R) string,
) ([]T, error) {
	results := make([]T, 0)
	nextPage := ""

	for {
		if nextPage != "" {
			query["page"] = nextPage
		}

		pageResp, err := getPageFunc(ctx, client, path, query)
		if err != nil {
			return nil, fmt.Errorf("failed to get page: %w", err)
		}

		results = append(results, getItems(pageResp)...)

		nextPage = getNextPage(pageResp)
		if nextPage == "" {
			break
		}
	}

	return results, nil
}

// Example of how to use paginateGeneric

// func ExamplePaginateUsersGeneric(ctx context.Context, client *Client) ([]User, error) {
// 	return paginateGeneric[User, paginatedResponse[User]](
// 		ctx,
// 		client,
// 		"/users",
// 		map[string]string{},
// 		pagedGetGeneric[paginatedResponse[User]],
// 		func(r *paginatedResponse[User]) []User { return r.Items },
// 		func(r *paginatedResponse[User]) string { return r.NextPage },
// 	)
// }
/*
type FooPage struct {
	Data []Foo `json:"data"`
	Next string `json:"next"`
}

page, _ := PagedGetGeneric[FooPage](ctx, client, "/foo", params)

items, _ := PaginateGeneric[Foo, FooPage](
	ctx,
	client,
	"/foo",
	params,
	PagedGetGeneric[FooPage],
	func(p *FooPage) []Foo { return p.Data },
	func(p *FooPage) string { return p.Next },
)

Those two functions are only needed for the generic paginator, because every external API returns pagination differently.

getItems tells the paginator how to extract the list of items from your custom response type.
getNextPage tells it how to get the “next page” token (or empty string to stop).

For the old PagedGet/Paginate, these are not needed because the response shape is fixed (items, next_page, total_count).
*/