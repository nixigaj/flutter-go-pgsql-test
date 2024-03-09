package main

import (
	"fmt"
	"net"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	PgConnStr     string `toml:"postgres_connect_string"`
	Bind          string `toml:"bind"`
	PgConnections int32  `toml:"postgres_parallel_connections"`
	Debug         bool   `toml:"debug"`
}

func getConfig(path string) (*Config, error) {
	var cfg Config

	configFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = toml.Unmarshal(configFile, &cfg)
	if err != nil {
		return nil, err
	}

	if cfg.Bind == "" {
		cfg.Bind = "localhost:8080"
	}

	if cfg.PgConnections < 0 {
		return nil, fmt.Errorf("negative number of paralell PostgreSQL connections: %d", cfg.PgConnections)
	}

	err = validateBind(cfg.Bind)
	if err != nil {
		return nil, fmt.Errorf("failed to validate bind: %v", err)
	}

	return &cfg, nil
}

func validateBind(bind string) error {
	listener, err := net.Listen("tcp", bind)
	if err != nil {
		return err
	}
	err = listener.Close()
	if err != nil {
		return err
	}
	return nil
}
