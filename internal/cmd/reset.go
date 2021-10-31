package cmd

import (
	"context"
	"github.com/aveloper/blog/internal/config"
	"github.com/aveloper/blog/internal/db"
	"github.com/aveloper/blog/internal/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
)

func reset(*cobra.Command, []string) {
	cfg := config.Get()
	log := logger.New(&logger.Config{Production: cfg.Production})

	dbCfg := &db.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		Name:     cfg.DB.Name,
		ForceTLS: false,
	}

	err := db.Down(context.Background(), dbCfg, log)
	if err != nil {
		log.Error("Failed to run the down migration", zap.Error(err))
		os.Exit(1)
	}

	log.Info("Successfully run the down migrations")
}