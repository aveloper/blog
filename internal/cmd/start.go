package cmd

import (
	"context"
	"os"

	"go.uber.org/zap"

	"github.com/aveloper/blog/internal/config"
	"github.com/aveloper/blog/internal/db"
	"github.com/aveloper/blog/internal/logger"
	"github.com/aveloper/blog/internal/server"
	"github.com/spf13/cobra"
)

func start(cmd *cobra.Command, _ []string) {
	apiOnly, err := cmd.Flags().GetBool("api-only")
	cobra.CheckErr(err)

	cfg := config.Get()
	log := logger.New(&logger.Config{Production: cfg.Production})

	dbConn, err := db.GetConnection(context.Background(), &db.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		Name:     cfg.DB.Name,
		ForceTLS: false,
	}, log)

	if err != nil {
		log.Error("failed to get DB connection", zap.Error(err))
		os.Exit(1)
	}

	srv := server.NewServer(&server.Config{
		Port:    cfg.Port,
		Logger:  log,
		APIOnly: apiOnly,
	}, dbConn)

	go srv.Listen()

	srv.WaitForShutdown()

	log.Info("Server Shutdown complete")
}
