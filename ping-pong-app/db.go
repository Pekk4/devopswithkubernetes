package main

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type CounterStore struct {
	db *sql.DB
}

func NewCounterStore(conn string) (*CounterStore, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &CounterStore{db: db}, nil
}

func (s *CounterStore) Close() error {
	return s.db.Close()
}

func (s *CounterStore) Init(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS pingpong_counter (
			id SERIAL PRIMARY KEY,
			value INT
		);
	`)
	return err
}

func (s *CounterStore) IncrementCounts(ctx context.Context) (int64, error) {
	var v int64
	err := s.db.QueryRowContext(ctx, "INSERT INTO pingpong_counter (value) VALUES (1) RETURNING id").Scan(&v)
	return v, err
}

func (s *CounterStore) GetCounts(ctx context.Context) (int64, error) {
	var v int64
	err := s.db.QueryRowContext(ctx, "SELECT id FROM pingpong_counter ORDER BY id DESC LIMIT 1").Scan(&v)
	return v, err
}

func (s *CounterStore) Ping() error {
	return s.db.Ping()
}
