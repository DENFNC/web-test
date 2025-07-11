package psql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Option func(cfg *databaseConfig)

type Database struct {
	*pgxpool.Pool
}

type databaseConfig struct {
	MaxConns          int32
	MinConns          int32
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
}

func WithMaxConns(max int32) Option {
	return func(cfg *databaseConfig) {
		cfg.MaxConns = max
	}
}

func WithMinConns(min int32) Option {
	return func(cfg *databaseConfig) {
		cfg.MinConns = min
	}
}

func WithMaxConnLifetime(t time.Duration) Option {
	return func(cfg *databaseConfig) {
		cfg.MaxConnLifetime = t
	}
}

func WithMaxConnIdleTime(t time.Duration) Option {
	return func(cfg *databaseConfig) {
		cfg.MaxConnIdleTime = t
	}
}

func WithHealthCheckPeriod(t time.Duration) Option {
	return func(cfg *databaseConfig) {
		cfg.HealthCheckPeriod = t
	}
}

func NewDatabase(ctx context.Context, connString string, options ...Option) (*Database, error) {
	var config databaseConfig
	for _, option := range options {
		option(&config)
	}

	configPool, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	configPool.MaxConns = config.MaxConns
	configPool.MinConns = config.MinConns
	configPool.MaxConnLifetime = config.MaxConnLifetime
	configPool.MaxConnIdleTime = config.MaxConnIdleTime
	configPool.HealthCheckPeriod = config.HealthCheckPeriod

	pool, err := pgxpool.NewWithConfig(ctx, configPool)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the database: %v", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &Database{
		Pool: pool,
	}, nil
}
