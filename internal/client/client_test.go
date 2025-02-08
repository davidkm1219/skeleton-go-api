package client_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/twk/skeleton-go-api/internal/client"
	mock_client "github.com/twk/skeleton-go-api/internal/client/mocks"
	"github.com/twk/skeleton-go-api/internal/logger"
)

func TestRequest(t *testing.T) {
	t.Parallel()

	type args struct {
		method    string
		targetURL string
		path      string
		headers   map[string]string
		query     map[string]string
		body      io.Reader
	}

	type fields struct {
		mockOps func(m *mock_client.MockhttpClient)
	}

	type wants struct {
		status int
		err    error
	}

	tests := map[string]struct {
		args   args
		fields fields
		wants  wants
	}{
		"Valid GET request": {
			args: args{
				method:    http.MethodGet,
				targetURL: "http://example.com/",
				path:      "/api/v1/resource",
				headers:   map[string]string{"Accept": "application/json"},
				query:     map[string]string{"id": "unique-id"},
				body:      io.NopCloser(strings.NewReader(`{"name": "snyk project"}`)),
			},
			fields: fields{
				mockOps: func(m *mock_client.MockhttpClient) {
					m.EXPECT().Do(&RequestMatcher{
						Method:   http.MethodGet,
						Host:     "example.com",
						Path:     "/api/v1/resource",
						RawQuery: "id=unique-id",
						Header:   http.Header{"Accept": []string{"application/json"}},
					}).Return(&http.Response{StatusCode: http.StatusOK}, nil)
				},
			},
			wants: wants{
				status: http.StatusOK,
				err:    nil,
			},
		},
		"Invalid URL": {
			args: args{
				method:    http.MethodGet,
				targetURL: "http://[::1]:namedport",
				path:      "/api/v1/resource",
			},
			fields: fields{
				mockOps: func(m *mock_client.MockhttpClient) {
					m.EXPECT().Do(gomock.Any()).Times(0)
				},
			},
			wants: wants{
				err: errors.New("failed to parse URL"),
			},
		},
		"Service error": {
			args: args{
				method:    http.MethodGet,
				targetURL: "http://example.com",
				path:      "/api/v1/resource",
				headers:   nil,
				query:     nil,
				body:      nil,
			},
			fields: fields{
				mockOps: func(m *mock_client.MockhttpClient) {
					m.EXPECT().Do(gomock.Any()).Return(nil, errors.New("service error"))
				},
			},
			wants: wants{
				status: 0,
				err:    errors.New("service error"),
			},
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockHTTPClient := mock_client.NewMockhttpClient(ctrl)
			tt.fields.mockOps(mockHTTPClient)

			log := logger.NewNop()
			c, err := client.NewClient(mockHTTPClient)
			assert.NoError(t, err)

			resp, err := c.Request(context.Background(), log, tt.args.method, tt.args.targetURL, tt.args.path, tt.args.headers, tt.args.query, tt.args.body)
			if tt.wants.err != nil {
				assert.ErrorContains(t, err, tt.wants.err.Error())
				return
			}

			assert.Equal(t, tt.wants.status, resp.StatusCode)
		})
	}
}

type RequestMatcher struct {
	Method   string
	Host     string
	Path     string
	RawQuery string
	Header   http.Header
}

func (r *RequestMatcher) Matches(x interface{}) bool {
	req, ok := x.(*http.Request)
	if !ok {
		return false
	}

	return r.matchMethod(req) && r.matchURL(req) && r.matchHeader(req)
}

func (r *RequestMatcher) matchMethod(req *http.Request) bool {
	return req.Method == r.Method
}

func (r *RequestMatcher) matchURL(req *http.Request) bool {
	return req.URL.Host == r.Host && req.URL.Path == r.Path && req.URL.RawQuery == r.RawQuery
}

func (r *RequestMatcher) matchHeader(req *http.Request) bool {
	for key, val := range r.Header {
		if req.Header.Get(key) != val[0] {
			return false
		}
	}

	return true
}

func (r *RequestMatcher) String() string {
	return fmt.Sprintf("RequestMatcher{Method: %s, Host: %s, Path: %s, RawQuery: %s, Header: %v}", r.Method, r.Host, r.Path, r.RawQuery, r.Header)
}
