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

//getConnConfig return a ConnConfig
func getConnConfig(cfg *Config, logger *zap.Logger) (*pgx.ConnConfig, error) {
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

	return connConfig, nil
}


//Get creates and returns a new DB connection
func Get(ctx context.Context, cfg *Config, logger *zap.Logger) (*DB, error) {

	connConfig, err := getConnConfig(cfg, logger)
	if err != nil {
		return nil, err
	}

	conn, err := pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		logger.Error("Failed to connect to DB", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	//execute the .up.sql files
	upMigrations(stdlib.OpenDB(*connConfig), logger)

	return &DB{Conn: conn}, nil
}

//Down to run the down migrations
func Down(ctx context.Context, cfg *Config, logger *zap.Logger) error {

	connConfig, err := getConnConfig(cfg, logger)
	if err != nil {
		return err
	}

	_, err = pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		logger.Error("Failed to connect to DB", zap.Error(err))
		return fmt.Errorf("failed to connect to DB: %w", err)
	}

	//execute the .down.sql files
	return downMigration(stdlib.OpenDB(*connConfig), logger)
}