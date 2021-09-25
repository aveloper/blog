package db

import (
	"blog/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/log/zapadapter"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type DB struct {
	*pgx.Conn
}

func Get(ctx context.Context, dbCfg *config.DB, logger *zap.Logger) (*DB, error) {
	// postgres://username:password@localhost:5432/database_name
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Name,
	)

	connConfig, err := pgx.ParseConfig(connString)
	if err != nil {
		logger.Error("Failed to parse connection string", zap.Error(err))
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	connConfig.Logger = zapadapter.NewLogger(logger)

	conn, err := pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		logger.Error("Failed to connect to DB", zap.Error(err))
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	return &DB{Conn: conn}, nil
}
