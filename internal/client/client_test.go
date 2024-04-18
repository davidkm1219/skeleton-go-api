package client_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twk/skeleton-go-api/internal/client"
)

func TestClient_Get(t *testing.T) {
	type fields struct {
		setup   func() *httptest.Server
		context func() context.Context
	}

	type want struct {
		resp *http.Response
		err  error
	}

	tests := map[string]struct {
		fields fields
		want   want
	}{
		"Successful GET request": {
			fields: fields{
				setup: func() *httptest.Server {
					return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						w.WriteHeader(http.StatusOK)
					}))
				},
				context: context.Background,
			},
			want: want{
				resp: &http.Response{
					StatusCode: http.StatusOK,
				},
			},
		},
		"Internal Server Error": {
			fields: fields{
				setup: func() *httptest.Server {
					return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						w.WriteHeader(http.StatusInternalServerError)
					}))
				},
				context: context.Background,
			},
			want: want{
				resp: &http.Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
		},
		"Request timeout": {
			fields: fields{
				setup: func() *httptest.Server {
					return httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
					}))
				},
				context: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					cancel()
					return ctx
				},
			},
			want: want{
				err: context.Canceled,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			server := tt.fields.setup()
			defer server.Close()

			httpClient := server.Client()
			c := client.NewClient(httpClient)

			resp, err := c.Get(tt.fields.context(), server.URL)
			if tt.want.err != nil {
				assert.ErrorContains(t, err, tt.want.err.Error())
				return
			}

			assert.Equal(t, tt.want.resp.StatusCode, resp.StatusCode)
		})
	}
}
