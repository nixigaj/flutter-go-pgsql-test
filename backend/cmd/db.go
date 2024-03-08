package main

import (
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
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

	err = initAndValidateDb(pool)
	if err != nil {
		return nil, fmt.Errorf("failed init and validate PostgreSQL: %v", err)
	}

	return pool, nil
}

func initAndValidateDb(pool *pgxpool.Pool) error {
	var exists bool
	err := pool.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'hello')").Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking table existence: %v", err)
	}

	if !exists {
		_, err := pool.Exec(context.Background(), "CREATE TABLE hello (id SERIAL PRIMARY KEY, counter INT NOT NULL, message TEXT NOT NULL)")
		if err != nil {
			return fmt.Errorf("error creating table: %v", err)
		}
	}

	var count int
	err = pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM hello").Scan(&count)
	if err != nil {
		if err != nil {
			return fmt.Errorf("error checking table size: %v", err)
		}
	}

	if count < 1 {
		_, err = pool.Exec(context.Background(), "INSERT INTO hello (counter, message) VALUES (0, 'Hello, world!')")
		if err != nil {
			if err != nil {
				return fmt.Errorf("error inserting: %v", err)
			}
		}
	}

	return nil
}

func incrementHello(pool *pgxpool.Pool) (string, error) {
	startTime := time.Now()

	_, err := pool.Exec(context.Background(), `UPDATE hello
	                                           SET counter = counter + 1
	                                           WHERE id = (SELECT MIN(id) FROM hello);`)
	if err != nil {
		return "", fmt.Errorf("exec fail: %v", err)
	}

	var id int
	var counter int
	var message string
	err = pool.QueryRow(context.Background(), "SELECT * FROM hello WHERE id = (SELECT MIN(id) FROM hello)").Scan(&id, &counter, &message)
	if err != nil {
		return "", fmt.Errorf("select fail: %v", err)
	}

	log.Debugf("Database increment counter latency: %dms", time.Since(startTime).Milliseconds())

	return fmt.Sprintf("ID: %d, Counter: %d, Message: %s", id, counter, message), nil
}
