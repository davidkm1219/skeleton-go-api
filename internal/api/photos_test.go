package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/twk/skeleton-go-api/internal/api"
	mock "github.com/twk/skeleton-go-api/internal/api/mocks"
	"github.com/twk/skeleton-go-api/internal/config"
	"github.com/twk/skeleton-go-api/internal/logger"
	"github.com/twk/skeleton-go-api/internal/photos"
)

func TestPhotosHandler(t *testing.T) {
	t.Parallel()

	type args struct {
		cfg *config.Server
		id  string
	}

	type fields struct {
		mockOperation func(m *mock.MockphotoService)
	}

	type want struct {
		code int
	}

	tests := map[string]struct {
		args   args
		fields fields
		want   want
	}{
		"success": {
			args: args{
				cfg: &config.Server{Timeout: 1 * time.Second},
				id:  "1",
			},
			fields: fields{
				mockOperation: func(m *mock.MockphotoService) {
					m.EXPECT().GetPhotos(gomock.Any(), 1).Return(&photos.Photo{}, nil)
				},
			},
			want: want{
				code: http.StatusOK,
			},
		},
		"invalid id": {
			args: args{
				cfg: &config.Server{Timeout: 1 * time.Second},
				id:  "abc",
			},
			fields: fields{
				mockOperation: func(m *mock.MockphotoService) {
					m.EXPECT().GetPhotos(gomock.Any(), 0).Times(0)
				},
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
		"service error": {
			args: args{
				cfg: &config.Server{Timeout: 1 * time.Second},
				id:  "1",
			},
			fields: fields{
				mockOperation: func(m *mock.MockphotoService) {
					m.EXPECT().GetPhotos(gomock.Any(), 1).Return(nil, assert.AnError)
				},
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mock.NewMockphotoService(ctrl)
			tt.fields.mockOperation(mockService)

			router := gin.Default()

			router.GET("/photos/:id", api.Photos(&config.Server{Timeout: 1 * time.Second}, mockService, logger.NewNop()))

			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/photos/"+tt.args.id, http.NoBody)
			assert.NoError(t, err)

			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)
			assert.Equal(t, tt.want.code, resp.Code)
		})
	}
}
