package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/log/zapadapter"

	"github.com/jackc/pgx/v4/stdlib"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type DB struct {
	*pgx.Conn
}

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	ForceTLS bool
}

//Get creates and returns a new DB connection
func Get(ctx context.Context, cfg *Config, logger *zap.Logger) (*DB, error) {
	sslMode := "disable"

	if cfg.ForceTLS {
		sslMode = "require"
	}

	// postgres://username:password@localhost:5432/database_name?sslmode=disable
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		sslMode,
	)

	connConfig, err := pgx.ParseConfig(connString)
	if err != nil {
		logger.Error("Failed to parse connection string", zap.Error(err))
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	connConfig.LogLevel = pgx.LogLevelDebug
	connConfig.Logger = zapadapter.NewLogger(logger)

	conn, err := pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		logger.Error("Failed to connect to DB", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	runMigrations(stdlib.OpenDB(*connConfig), logger)

	return &DB{Conn: conn}, nil
}