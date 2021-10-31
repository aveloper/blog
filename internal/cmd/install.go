package cmd

import (
	"os"

	"go.uber.org/zap"

	"github.com/aveloper/blog/internal/config"
	"github.com/aveloper/blog/internal/db"
	"github.com/aveloper/blog/internal/logger"

	"github.com/spf13/cobra"
)

func install(cmd *cobra.Command, _ []string) {
	// TODO: Figure out how force will work and implement
	force, err := cmd.Flags().GetBool("force")
	cobra.CheckErr(err)

	// If the command install, was called with force flag
	// reset the config
	if force {
		config.Reset()
	}

	// Install steps
	// 1. Read config and create if not already present
	// 2. Connect to DB and run migrations
	// 3. Setup autostart
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

	err = db.Setup(dbCfg, log)
	if err != nil {
		log.Error("DB Error", zap.Error(err))
		os.Exit(1)
	}

	// TODO: Setup autostart
}
