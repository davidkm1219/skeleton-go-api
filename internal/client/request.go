package client

//go:generate mockgen -destination=mocks/request.go -package=mock_client -source=request.go

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/twk/skeleton-go-api/internal/logger"
)

const basicAuthCredentialParts = 2

// HTTPRequester is an interface for making HTTP requests.
type HTTPRequester interface {
	Request(ctx context.Context, logger *logger.Logger, method, url, path string, header, query map[string]string, body io.Reader) (*http.Response, error)
}

// Get makes a GET request to the target URL with the specified query parameters and returns the response body.
func Get[T any](ctx context.Context, log *logger.Logger, c HTTPRequester, targetURL, path string, query map[string]string, authType AuthType, credential *string) (resp *T, code int, err error) {
	header := map[string]string{
		"Accept": "application/json",
	}

	if aErr := setAuth(authType, credential, header, log); aErr != nil {
		return nil, 0, fmt.Errorf("failed to set auth: %w", err)
	}

	r, err := c.Request(ctx, log, http.MethodGet, targetURL, path, header, query, nil)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to send request: %w", err)
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return nil, r.StatusCode, fmt.Errorf("unexpected status code: %d", r.StatusCode)
	}

	responseBody, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to read response body: %w", err)
	}

	var res T

	if err := json.Unmarshal(responseBody, &res); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &res, r.StatusCode, nil
}

// Post makes a POST request to the target URL with the specified query parameters and body and returns the response body.
func Post[B any, T any](ctx context.Context, log *logger.Logger, c HTTPRequester, targetURL, path string, query map[string]string, body *B, authType AuthType, credential *string) (resp *T, code int, err error) {
	header := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}

	if aErr := setAuth(authType, credential, header, log); aErr != nil {
		return nil, 0, fmt.Errorf("failed to set auth: %w", aErr)
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	r, err := c.Request(ctx, log, http.MethodPost, targetURL, path, header, query, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to send request: %w", err)
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return nil, r.StatusCode, fmt.Errorf("unexpected status code: %d", r.StatusCode)
	}

	responseBody, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to read response body: %w", err)
	}

	var res T

	if err := json.Unmarshal(responseBody, &res); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return &res, r.StatusCode, nil

}

func setAuth(authType AuthType, credential *string, header map[string]string, log *logger.Logger) error {
	if credential == nil {
		log.Info("No credential provided")
		return nil
	}

	switch authType {
	case AuthTypeToken:
		header["Authorization"] = "Token " + *credential
	case AuthTypeBearer:
		header["Authorization"] = "Bearer " + *credential
	case AuthTypeBasic:
		cred, err := parseBasicAuth(*credential)
		if err != nil {
			return fmt.Errorf("failed to parse basic auth: %w", err)
		}

		header["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(cred))
	}

	return nil
}

func parseBasicAuth(credential string) (string, error) {
	creds := strings.Split(credential, ":")
	if len(creds) != basicAuthCredentialParts {
		return "", fmt.Errorf("invalid basic auth credential")
	}

	return credential, nil
}
