// Package api provides the handlers for the API endpoints.
package api

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/twk/skeleton-go-api/internal/config"
	"github.com/twk/skeleton-go-api/internal/photos"
)

type photoService interface {
	GetPhotos(ctx context.Context, albumID int) (*photos.Photo, error)
}

// Photos returns a handler for getting photos
func Photos(cfg *config.Server, ps photoService, l *zap.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), cfg.Timeout)
		defer cancel()

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			l.Error("failed to parse id", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})

			return
		}

		p, err := ps.GetPhotos(ctx, id)
		if err != nil {
			l.Error("failed to get photos", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get photos"})

			return
		}

		c.JSON(http.StatusOK, p)
	}
}
