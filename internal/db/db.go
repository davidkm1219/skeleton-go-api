package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/universe/thrubble-api-go/internal/config"
)

// DatabasePool holds the database connection pool.
type DatabasePool struct {
	DB *sqlx.DB
}

// NewDatabasePool creates a new database connection pool.
func NewDatabasePool(cfg *config.Config) (*DatabasePool, error) {
	db, err := sql.Open("postgres", cfg.Database.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	dbx := sqlx.NewDb(db, "postgres")

	applyPoolSettings(db, cfg)
	if err := pingWithTimeout(db, cfg.Database.PingTimeout); err != nil {
		return nil, err
	}

	return &DatabasePool{DB: dbx}, nil
}

// Close shuts down the database pool.
func (p *DatabasePool) Close() error {
	if p == nil || p.DB == nil {
		return nil
	}
	return p.DB.Close()
}

func applyPoolSettings(db *sql.DB, cfg *config.Config) {
	if cfg.Database.MaxConnection > 0 {
		db.SetMaxOpenConns(cfg.Database.MaxConnection)
	}
	if cfg.Database.MaxIdleConnection > 0 {
		db.SetMaxIdleConns(cfg.Database.MaxIdleConnection)
	}
	if cfg.Database.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)
	}
}

func pingWithTimeout(db *sql.DB, timeout time.Duration) error {
	if timeout <= 0 {
		return db.Ping()
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return db.PingContext(ctx)
}
