package db

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"
)

// migrationVersion defines the current migration version. This ensures the app
// is always compatible with the version of the database.
const migrationVersion = 1

//go:embed migrations/*.sql
var migrations embed.FS

//getNewMigration get the new migrate instance 
func getNewMigration(db *sql.DB) (source.Driver, database.Driver, *migrate.Migrate, error) {
	// Read the source for the migrations.
	// Our source is the SQL files in the migrations folder
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read migration source %w",  err)
	}

	// Connect with the target i.e, our postgres DB
	target, err := postgres.WithInstance(db, new(postgres.Config))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to read migration target %w", err)
	}

	// Create a new instance of the migration using the defined source and target
	m, err := migrate.NewWithInstance("iofs", source, "postgres", target)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create migration instance %w", err)
	}

	return source, target, m, nil
}

//upMigrations migrates the postgres schema to the current version
func upMigrations(db *sql.DB, logger *zap.Logger) {

	source, target, m, err := getNewMigration(db)
	if  err != nil {
		logger.Panic(" failed to run the migration ",zap.Error(err))
	}
	// Ensure that we close the source connection
	defer func() {
		err := source.Close()
		if err != nil {
			logger.Error("failed to close source connection", zap.Error(err))
		}
	}()

	// Ensure that we close the target connection
	defer func() {
		err := target.Close()
		if err != nil {
			logger.Error("failed to close target connection", zap.Error(err))
		}
	}()

	// Migrate the DB to the current version
	err = m.Migrate(migrationVersion)
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		//If the error is present, and the error is not the No change error
		// log the error
		logger.Panic("failed to run migration", zap.Error(err), zap.Int("version", migrationVersion))
	}

	if err != nil && errors.Is(err, migrate.ErrNoChange) {
		logger.Info("Database is already migrated. Nothing to change", zap.Int("version", migrationVersion))
		return
	}

	logger.Info("Successfully executed migrations for DB", zap.Int("version", migrationVersion))
}

//downMigration migrates the postgres schema from the active migration version to all the way down
func downMigration(db *sql.DB, logger *zap.Logger) error {
	
	source, target, m, err := getNewMigration(db)
	if  err != nil {
		logger.Panic(" failed to run the down migration ",zap.Error(err))
	}

	// Ensure that we close the source connection
	defer func() {
		err := source.Close()
		if err != nil {
			logger.Error("failed to close source connection", zap.Error(err))
		}
	}()

	// Ensure that we close the target connection
	defer func() {
		err := target.Close()
		if err != nil {
			logger.Error("failed to close target connection", zap.Error(err))
		}
	}()
	
	//applying all down migrations
	err = m.Down()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		//If the error is present, and the error is not the No change error
		// log the error
		logger.Panic("failed to run migration", zap.Error(err), zap.Int("version", migrationVersion))
	}

	if err != nil && errors.Is(err, migrate.ErrNoChange) {
		logger.Info("Database is already down migrated. Nothing to change", zap.Int("version", migrationVersion))
	}

	return nil
}