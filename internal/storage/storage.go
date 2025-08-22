package storage

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewPool(ctx context.Context, dsn string) *Storage {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		// If the DSN is invalid, log the error and exit
		log.Fatalf("Unable to parse database URL: %v\n", err)
	}

	cfg.MaxConns = 10 // Set a maximum number of connections
	cfg.MinConns = 1  // Set a minimum number of connections
	cfg.HealthCheckPeriod = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	return &Storage{pool: pool}
}

func (s *Storage) Close() {
	if s.pool != nil {
		s.pool.Close()
	}
}
