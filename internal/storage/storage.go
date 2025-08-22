package storage

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/agidelle/merch-shop/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewPool(ctx context.Context, dsn string) *Storage {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Unable to parse database URL: %v\n", err)
	}

	cfg.MaxConns = 10 // Set a maximum number of connections
	cfg.MinConns = 1  // Set a minimum number of connections
	cfg.HealthCheckPeriod = 30 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	// Check if the connection pool is working
	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	return &Storage{pool: pool}
}

func (s *Storage) Close() {
	if s.pool != nil {
		s.pool.Close()
	}
}

func (s *Storage) FindTask(filter *domain.Filter) ([]*domain.Task, error) {
	tasks := make([]*domain.Task, 0)
	query := "SELECT id, date, title, comment, repeat FROM scheduler"
	args := []interface{}{}
	conditions := []string{}

	//Добавление условий в зависимости от фильтра
	if filter.ID != nil {
		conditions = append(conditions, "id = ?")
		args = append(args, *filter.ID)
	}
	if filter.SearchTerm != "" {
		searchPattern := "%" + filter.SearchTerm + "%"
		conditions = append(conditions, "(title LIKE ? OR comment LIKE ?)")
		args = append(args, searchPattern, searchPattern)
	}
	if filter.Date != "" {
		conditions = append(conditions, "date = ?")
		args = append(args, filter.Date)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY date"
	if filter.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, filter.Limit)
	}
}
