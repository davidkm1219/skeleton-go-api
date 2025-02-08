package photos

import (
	"context"
	"io"
	"net/http"
	"strconv"

	hClient "github.com/twk/skeleton-go-api/internal/client"
	"github.com/twk/skeleton-go-api/internal/logger"
)

const (
	PhotoBaseURL = "https://jsonplaceholder.typicode.com"
	photoPath    = "/photos"
)

type httpClient interface {
	Request(ctx context.Context, logger *logger.Logger, method, url, path string, header, query map[string]string, body io.Reader) (*http.Response, error)
}

// PhotoClient is a client for the photo API.
type PhotoClient struct {
	baseURL    string
	authType   hClient.AuthType
	httpClient httpClient
	log        *logger.Logger
}

// NewClient creates a new photo client.
func NewClient(baseURL string, authType hClient.AuthType, httpClient httpClient, log *logger.Logger) *PhotoClient {
	return &PhotoClient{
		baseURL:    baseURL,
		authType:   authType,
		httpClient: httpClient,
		log:        log,
	}
}

// GetPhotos gets photos from the API.
func (c *PhotoClient) GetPhotos(ctx context.Context, id int) (*Photo, error) {
	query := map[string]string{
		"albumId": strconv.Itoa(id),
	}

	photo, _, err := hClient.Get[Photo](ctx, c.log, c.httpClient, c.baseURL, photoPath, query, c.authType, nil)
	if err != nil {
		return nil, err
	}

	return photo, nil
}
