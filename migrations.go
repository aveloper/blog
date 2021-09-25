package main

import (
	"blog/db"
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v4/stdlib"

	"go.uber.org/zap"
)

// migrationVersion defines the current migration version. This ensures the app
// is always compatible with the version of the database.
const migrationVersion = 1

//go:embed migrations/*.sql
var migrations embed.FS

//runMigrations migrates the postgres schema to the current version
func runMigrations(db *db.DB, logger *zap.Logger) {
	// Read the source for the migrations.
	// Our source is the SQL files in the migrations folder
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		logger.Panic("failed to read migration source", zap.Error(err))
		return
	}

	// Ensure that we close the source connection
	defer func() {
		err := source.Close()
		if err != nil {
			logger.Error("failed to close source connection", zap.Error(err))
		}
	}()

	// Connect with the target i.e, our postgres DB
	target, err := postgres.WithInstance(stdlib.OpenDB(*db.Config()), new(postgres.Config))
	if err != nil {
		logger.Panic("failed to read migration target", zap.Error(err))
		return
	}

	// Ensure that we clos the target connection
	defer func() {
		err := target.Close()
		if err != nil {
			logger.Error("failed to close target connection", zap.Error(err))
		}
	}()

	// Create a new instance of the migration using the defined source and target
	m, err := migrate.NewWithInstance("iofs", source, "postgres", target)
	if err != nil {
		logger.Panic("failed to create migration instance", zap.Error(err))
		return
	}

	// Migrate the DB to the current version
	err = m.Migrate(migrationVersion)
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		//If the error is present, and the error is not the No change error
		// log the error
		logger.Panic("failed to run migration", zap.Error(err), zap.Int("version", migrationVersion))
	}

	logger.Info("Successfully executed migrations", zap.Int("version", migrationVersion))
}
