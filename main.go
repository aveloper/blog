package main

import (
	"blog/config"
	"blog/db"
	"context"
	"log"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	// Read the config
	appConfig := config.Get()

	// Build logger
	logger := getLogger(appConfig.Production, appConfig.Logger)

	// Defer flushing the logger
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			log.Printf("Failed to buffer zap logger: %v", err)
		}
	}(logger)

	// Connect to DB
	dbConn, err := db.Get(ctx, appConfig.DB, logger)
	if err != nil {
		logger.Fatal("DB error", zap.Error(err))
	}

	// Defer closing the DB
	defer func(dbConn *db.DB, ctx context.Context) {
		err := dbConn.Close(ctx)
		if err != nil {
			logger.Fatal("failed to close DB", zap.Error(err))
		}
	}(dbConn, ctx)

	logger.Info("Successfully connected to DB")

	// Execute the migrations
	runMigrations(dbConn, logger)

	// Build the server
	s := NewServer(appConfig, logger)
	s.Initialize()

	// Start the server
	s.Listen()

	// Wait for the server to close
	<-s.connClose
}
