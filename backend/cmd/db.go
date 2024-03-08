package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func startDbPool(cfg *Config) (*pgxpool.Pool, error) {
	pgCfg, err := pgxpool.ParseConfig(cfg.PgConnStr)
	if err != nil {
		return nil, fmt.Errorf("PostgreSQL config parsing error: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgCfg)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to PostgreSQL: %v", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %v", err)
	}

	return pool, nil
}
