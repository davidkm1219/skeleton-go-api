// Package photos provides the operations for handling photos operations. It contains the Service struct and the GetPhotosConcurrently function.
package photos

//go:generate mockgen -destination=mocks/photos.go -package=mock_photos -source=photos.go

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"

	"github.com/twk/skeleton-go-api/internal/logger"
)

// Photo represents a photo object
type Photo struct {
	AlbumID      int    `json:"albumId"`
	ID           int    `json:"id"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

// Result represents the result of a photo operation
type Result struct {
	Photo *Photo
	Err   error
}

type client interface {
	GetPhotos(ctx context.Context, id int) (*Photo, error)
}

// Service provides the operations for handling photos operations
type Service struct {
	client client
	log    *logger.Logger
}

// NewService creates a new Service for handling photos operations
func NewService(c client, log *logger.Logger) *Service {
	return &Service{
		client: c,
		log:    log,
	}
}

// GetPhotosConcurrently gets photos concurrently
func (s *Service) GetPhotosConcurrently(ctx context.Context, concurrency int) []int {
	var wg sync.WaitGroup

	wg.Add(concurrency)

	chanResult := make(chan Result)

	for i := 1; i <= concurrency; i++ {
		go func(id int) {
			defer wg.Done()

			photo, err := s.client.GetPhotos(ctx, id)
			chanResult <- Result{Photo: photo, Err: err}
		}(i)
	}

	// We can just use buffered channel and wg.Wait() here. Depends on how we want to handle the result
	// Doing this way, so we are processing the result as soon as it is available
	go func() {
		wg.Wait()
		close(chanResult)
	}()

	processedPhotos := make([]int, 0)

	for r := range chanResult {
		if r.Err != nil {
			s.log.Error("Failed to process photo", zap.Error(r.Err))
			continue
		}

		s.log.Info("Processed photo", zap.Int("id", r.Photo.ID))
		processedPhotos = append(processedPhotos, r.Photo.ID)
	}

	return processedPhotos
}

// GetPhotos gets photos from the photos URL
func (s *Service) GetPhotos(ctx context.Context, id int) (*Photo, error) {
	resp, err := s.client.GetPhotos(ctx, id)
	if err != nil {
		s.log.Error("Failed to get photos", zap.Error(err))
		return nil, fmt.Errorf("failed to get photos: %w", err)
	}

	return resp, nil
}
