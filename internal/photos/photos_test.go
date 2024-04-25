package photos_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/twk/skeleton-go-api/internal/logger"
	"github.com/twk/skeleton-go-api/internal/photos"
	mock_photos "github.com/twk/skeleton-go-api/internal/photos/mocks"
)

func TestGetPhotos(t *testing.T) {
	type fields struct {
		mockOperation func(m *mock_photos.Mockclient)
	}

	type want struct {
		want *photos.Photo
		err  error
	}

	tests := map[string]struct {
		fields fields
		want   want
	}{
		"success": {
			fields: fields{
				mockOperation: func(m *mock_photos.Mockclient) {
					m.EXPECT().Get(context.Background(), "https://jsonplaceholder.typicode.com/photos/1").Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(`{"albumId":1,"id":1,"title":"test","url":"test","thumbnailUrl":"test"}`))),
					}, nil)
				},
			},
			want: want{want: &photos.Photo{AlbumID: 1, ID: 1, Title: "test", URL: "test", ThumbnailURL: "test"}},
		},
		"error": {
			fields: fields{
				mockOperation: func(m *mock_photos.Mockclient) {
					m.EXPECT().Get(context.Background(), "https://jsonplaceholder.typicode.com/photos/1").Return(nil, errors.New("error"))
				},
			},
			want: want{err: errors.New("failed to get photos: error")},
		},
		"http not OK": {
			fields: fields{
				mockOperation: func(m *mock_photos.Mockclient) {
					m.EXPECT().Get(context.Background(), "https://jsonplaceholder.typicode.com/photos/1").Return(&http.Response{
						StatusCode: http.StatusNotFound,
						Body:       io.NopCloser(bytes.NewReader([]byte(``))),
					}, nil)
				},
			},
			want: want{err: errors.New("received non-OK HTTP status: 404")},
		},
		"invalid body": {
			fields: fields{
				mockOperation: func(m *mock_photos.Mockclient) {
					m.EXPECT().Get(context.Background(), "https://jsonplaceholder.typicode.com/photos/1").Return(&http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(`{"albumId":1,"id":1,"title":"test","url":"test","thumbnailUrl":}`))),
					}, nil)
				},
			},
			want: want{err: errors.New("failed to decode response body: invalid character '}' looking for beginning of value")},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cl := mock_photos.NewMockclient(ctrl)
			tt.fields.mockOperation(cl)

			s := photos.NewService(cl, logger.NewNop())

			result, err := s.GetPhotos(context.Background(), 1)
			if tt.want.err != nil {
				assert.EqualError(t, err, tt.want.err.Error())
				return
			}

			assert.Equal(t, tt.want.want, result)
		})
	}
}

func TestGetPhotosConcurrently(t *testing.T) {
	type args struct {
		concurrency int
	}

	type fields struct {
		mockOperation func(m *mock_photos.Mockclient)
	}

	type want struct {
		want []int
	}

	tests := map[string]struct {
		args   args
		fields fields
		want   want
	}{
		"success": {
			args: args{concurrency: 5},
			fields: fields{
				mockOperation: func(m *mock_photos.Mockclient) {
					for i := 1; i <= 5; i++ {
						m.EXPECT().Get(context.Background(), fmt.Sprintf("https://jsonplaceholder.typicode.com/photos/%d", i)).Return(&http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewReader([]byte(fmt.Sprintf(`{"albumId":1,"id":%d,"title":"test","url":"test","thumbnailUrl":"test"}`, i)))),
						}, nil)
					}
				},
			},
			want: want{want: []int{1, 2, 3, 4, 5}},
		},
		"error": {
			args: args{concurrency: 5},
			fields: fields{
				mockOperation: func(m *mock_photos.Mockclient) {
					m.EXPECT().Get(context.Background(), "https://jsonplaceholder.typicode.com/photos/1").Return(nil, errors.New("error"))
					for i := 2; i <= 5; i++ {
						m.EXPECT().Get(context.Background(), fmt.Sprintf("https://jsonplaceholder.typicode.com/photos/%d", i)).Return(&http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewReader([]byte(fmt.Sprintf(`{"albumId":1,"id":%d,"title":"test","url":"test","thumbnailUrl":"test"}`, i)))),
						}, nil)
					}
				},
			},
			want: want{want: []int{2, 3, 4, 5}},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cl := mock_photos.NewMockclient(ctrl)
			tt.fields.mockOperation(cl)

			s := photos.NewService(cl, logger.NewNop())

			result := s.GetPhotosConcurrently(context.Background(), tt.args.concurrency)

			assert.ElementsMatch(t, tt.want.want, result)
		})
	}
}
