package client_test

import (
	"context"
	"errors"
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

func TestGet(t *testing.T) {
	type args struct {
		targetURL  string
		path       string
		query      map[string]string
		authType   client.AuthType
		credential *string
	}

	type fields struct {
		mockOps func(m *mock_client.MockHTTPRequester)
	}

	type wants struct {
		resp *map[string]string
		code int
		err  error
	}

	tests := map[string]struct {
		args   args
		fields fields
		wants  wants
	}{
		"Successful GET request": {
			args: args{
				targetURL: "http://example.com",
				path:      "/api/v1/resource",
				query:     map[string]string{"key": "value"},
				authType:  client.AuthTypeBearer,
				credential: func() *string {
					s := "token"
					return &s
				}(),
			},
			fields: fields{
				mockOps: func(m *mock_client.MockHTTPRequester) {
					m.EXPECT().Request(gomock.Any(), gomock.Any(), http.MethodGet, "http://example.com", "/api/v1/resource", gomock.Any(), map[string]string{"key": "value"}, nil).
						Return(&http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(strings.NewReader(`{"name": "example"}`)),
						}, nil)
				},
			},
			wants: wants{
				resp: &map[string]string{"name": "example"},
				code: http.StatusOK,
				err:  nil,
			},
		},
		"Invalid URL": {
			args: args{
				targetURL: "http://[::1]:namedport",
				path:      "/api/v1/resource",
			},
			fields: fields{
				mockOps: func(m *mock_client.MockHTTPRequester) {
					m.EXPECT().Request(gomock.Any(), gomock.Any(), http.MethodGet, "http://[::1]:namedport", "/api/v1/resource", gomock.Any(), gomock.Any(), nil).
						Return(nil, errors.New("failed to parse URL"))
				},
			},
			wants: wants{
				resp: nil,
				code: http.StatusInternalServerError,
				err:  errors.New("failed to send request: failed to parse URL"),
			},
		},
		"Service error": {
			args: args{
				targetURL: "http://example.com",
				path:      "/api/v1/resource",
			},
			fields: fields{
				mockOps: func(m *mock_client.MockHTTPRequester) {
					m.EXPECT().Request(gomock.Any(), gomock.Any(), http.MethodGet, "http://example.com", "/api/v1/resource", gomock.Any(), gomock.Any(), nil).
						Return(nil, errors.New("service error"))
				},
			},
			wants: wants{
				resp: nil,
				code: http.StatusInternalServerError,
				err:  errors.New("failed to send request: service error"),
			},
		},
		"Unexpected status code": {
			args: args{
				targetURL: "http://example.com",
				path:      "/api/v1/resource",
			},
			fields: fields{
				mockOps: func(m *mock_client.MockHTTPRequester) {
					m.EXPECT().Request(gomock.Any(), gomock.Any(), http.MethodGet, "http://example.com", "/api/v1/resource", gomock.Any(), gomock.Any(), nil).
						Return(&http.Response{
							StatusCode: http.StatusBadRequest,
							Body:       io.NopCloser(strings.NewReader(``)),
						}, nil)
				},
			},
			wants: wants{
				resp: nil,
				code: http.StatusBadRequest,
				err:  errors.New("unexpected status code: 400"),
			},
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRequester := mock_client.NewMockHTTPRequester(ctrl)
			tt.fields.mockOps(mockRequester)

			log := logger.NewNop()

			resp, code, err := client.Get[map[string]string](context.Background(), log, mockRequester, tt.args.targetURL, tt.args.path, tt.args.query, tt.args.authType, tt.args.credential)
			if tt.wants.err != nil {
				assert.ErrorContains(t, err, tt.wants.err.Error())
				return
			}

			assert.NoError(t, err)

			assert.Equal(t, tt.wants.code, code)
			assert.Equal(t, tt.wants.resp, resp)
		})
	}
}

func TestPost(t *testing.T) {
	type args struct {
		targetURL  string
		path       string
		query      map[string]string
		body       *map[string]string
		authType   client.AuthType
		credential *string
	}

	type fields struct {
		mockOps func(m *mock_client.MockHTTPRequester)
	}

	type wants struct {
		resp *map[string]string
		code int
		err  error
	}

	tests := map[string]struct {
		args   args
		fields fields
		wants  wants
	}{
		"Successful POST request": {
			args: args{
				targetURL: "http://example.com",
				path:      "/api/v1/resource",
				query:     map[string]string{"key": "value"},
				body: &map[string]string{
					"name": "example",
				},
				authType: client.AuthTypeBearer,
				credential: func() *string {
					s := "token"
					return &s
				}(),
			},
			fields: fields{
				mockOps: func(m *mock_client.MockHTTPRequester) {
					m.EXPECT().Request(gomock.Any(), gomock.Any(), http.MethodPost, "http://example.com", "/api/v1/resource", gomock.Any(), map[string]string{"key": "value"}, gomock.Any()).
						Return(&http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(strings.NewReader(`{"name": "example"}`)),
						}, nil)
				},
			},
			wants: wants{
				resp: &map[string]string{"name": "example"},
				code: http.StatusOK,
				err:  nil,
			},
		},
		"Invalid URL": {
			args: args{
				targetURL: "http://[::1]:namedport",
				path:      "/api/v1/resource",
				body:      &map[string]string{"name": "example"},
			},
			fields: fields{
				mockOps: func(m *mock_client.MockHTTPRequester) {
					m.EXPECT().Request(gomock.Any(), gomock.Any(), http.MethodPost, "http://[::1]:namedport", "/api/v1/resource", gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil, errors.New("failed to parse URL"))
				},
			},
			wants: wants{
				resp: nil,
				code: http.StatusInternalServerError,
				err:  errors.New("failed to send request: failed to parse URL"),
			},
		},
		"Service error": {
			args: args{
				targetURL: "http://example.com",
				path:      "/api/v1/resource",
				body:      &map[string]string{"name": "example"},
			},
			fields: fields{
				mockOps: func(m *mock_client.MockHTTPRequester) {
					m.EXPECT().Request(gomock.Any(), gomock.Any(), http.MethodPost, "http://example.com", "/api/v1/resource", gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil, errors.New("service error"))
				},
			},
			wants: wants{
				resp: nil,
				code: http.StatusInternalServerError,
				err:  errors.New("failed to send request: service error"),
			},
		},
		"Unexpected status code": {
			args: args{
				targetURL: "http://example.com",
				path:      "/api/v1/resource",
				body:      &map[string]string{"name": "example"},
			},
			fields: fields{
				mockOps: func(m *mock_client.MockHTTPRequester) {
					m.EXPECT().Request(gomock.Any(), gomock.Any(), http.MethodPost, "http://example.com", "/api/v1/resource", gomock.Any(), gomock.Any(), gomock.Any()).
						Return(&http.Response{
							StatusCode: http.StatusBadRequest,
							Body:       io.NopCloser(strings.NewReader(``)),
						}, nil)
				},
			},
			wants: wants{
				resp: nil,
				code: http.StatusBadRequest,
				err:  errors.New("unexpected status code: 400"),
			},
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRequester := mock_client.NewMockHTTPRequester(ctrl)
			tt.fields.mockOps(mockRequester)

			log := logger.NewNop()

			resp, code, err := client.Post[map[string]string, map[string]string](context.Background(), log, mockRequester, tt.args.targetURL, tt.args.path, tt.args.query, tt.args.body, tt.args.authType, tt.args.credential)
			if tt.wants.err != nil {
				assert.ErrorContains(t, err, tt.wants.err.Error())
				return
			}

			assert.NoError(t, err)

			assert.Equal(t, tt.wants.code, code)
			assert.Equal(t, tt.wants.resp, resp)
		})
	}
}
