package db

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/stdlib"

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

//Setup runs the database migrations upto the current version
func Setup(cfg *Config, logger *zap.Logger) error {
	pgxConfig, err := getPgxConfig(cfg, logger)
	if err != nil {
		return err
	}

	runMigrations(stdlib.OpenDB(*pgxConfig), logger)

	return nil
}

//NewConnection creates and returns a new DB connection
func NewConnection(ctx context.Context, cfg *Config, logger *zap.Logger) (*DB, error) {
	return connect(ctx, cfg, logger)
}

var (
	once   sync.Once
	dbConn *DB
)

//GetConnection creates a new connection if not available or else returns the current connection
func GetConnection(ctx context.Context, cfg *Config, logger *zap.Logger) (*DB, error) {
	var (
		conn *DB
		err  error
	)

	once.Do(func() {
		conn, err = NewConnection(ctx, cfg, logger)
	})

	dbConn = conn

	return dbConn, nil
}

//getPgxConfig builds and returns the pgx connection config
func getPgxConfig(cfg *Config, logger *zap.Logger) (*pgx.ConnConfig, error) {
	sslMode := "prefer"

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

	// FIXME: The log level is always INFO, even when we set to Debug here
	connConfig.LogLevel = pgx.LogLevelDebug
	connConfig.Logger = zapadapter.NewLogger(logger)

	return connConfig, nil
}

//connect creates a new DB connection
func connect(ctx context.Context, cfg *Config, logger *zap.Logger) (*DB, error) {
	pgxConfig, err := getPgxConfig(cfg, logger)
	if err != nil {
		return nil, err
	}

	conn, err := pgx.ConnectConfig(ctx, pgxConfig)
	if err != nil {
		logger.Error("Failed to connect to DB", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	return &DB{
		Conn: conn,
	}, nil
}
