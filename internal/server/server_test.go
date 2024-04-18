package server_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/twk/skeleton-go-api/internal/config"
	"github.com/twk/skeleton-go-api/internal/server"
)

func TestServerServeHTTP(t *testing.T) {
	t.Parallel()

	type args struct {
		method string
		path   string
	}

	type want struct {
		status int
	}

	tests := map[string]struct {
		args args
		want want
	}{
		"RootPath": {args: args{method: http.MethodGet, path: "/"}, want: want{status: http.StatusOK}},
		"NotFound": {args: args{method: http.MethodGet, path: "/notfound"}, want: want{status: http.StatusNotFound}},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			logger := zap.NewNop()
			router := gin.Default()
			s := server.NewServer(&config.Server{Port: 8080}, router, []server.RouteParam{}, logger)

			req, err := http.NewRequestWithContext(context.Background(), tt.args.method, tt.args.path, http.NoBody)
			assert.NoError(t, err)

			resp := httptest.NewRecorder()

			s.ServeHTTP(resp, req)

			assert.Equal(t, tt.want.status, resp.Code)
		})
	}
}

func TestLoggerMiddleware(t *testing.T) {
	logger := zap.NewNop()
	router := gin.Default()
	s := server.NewServer(&config.Server{Port: 8080}, router, []server.RouteParam{}, logger)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", http.NoBody)
	assert.NoError(t, err)

	resp := httptest.NewRecorder()

	router.Use(s.LoggerMiddleware())
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}
