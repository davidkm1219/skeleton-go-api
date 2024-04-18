// Package server provides the HTTP server for the application. It contains the Server struct and the NewServer function.
package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/twk/skeleton-go-api/internal/config"
)

// RouteParam holds the each service that is required for the routes.
type RouteParam struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
}

type httpRouter interface {
	GET(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	POST(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	PUT(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	DELETE(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	NoRoute(handlers ...gin.HandlerFunc)
	Use(middleware ...gin.HandlerFunc) gin.IRoutes
	Run(addr ...string) error
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

// Server represents the HTTP server.
type Server struct {
	config *config.Server
	router httpRouter
	log    *zap.Logger
}

// NewServer creates a new server instance.
func NewServer(cfg *config.Server, r httpRouter, rp []RouteParam, log *zap.Logger) *Server {
	server := &Server{
		config: cfg,
		router: r,
		log:    log,
	}
	server.registerMiddleware()
	server.registerRoutes(rp)

	return server
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	err := s.router.Run(fmt.Sprintf("%s:%d", s.config.Host, s.config.Port))
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) registerRoutes(rp []RouteParam) {
	s.router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	for _, r := range rp {
		switch r.Method {
		case http.MethodGet:
			s.router.GET(r.Path, r.Handler)
		case http.MethodPost:
			s.router.POST(r.Path, r.Handler)
		case http.MethodPut:
			s.router.PUT(r.Path, r.Handler)
		case http.MethodDelete:
			s.router.DELETE(r.Path, r.Handler)
		}
	}

	s.router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
	})

	// Register middlewares
	s.router.Use(s.LoggerMiddleware())
}

func (s *Server) registerMiddleware() {
	s.router.Use(s.LoggerMiddleware())
}

// LoggerMiddleware instances a Logger middleware for Gin.
func (s *Server) LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = fmt.Sprintf("%s?%s", path, raw)
		}

		s.log.Debug("http request", zap.String("method", method), zap.String("path", path), zap.Int("status", statusCode), zap.Duration("latency", latency))
	}
}
