package cmd

import (
	"github.com/aveloper/blog/internal/config"
	"github.com/aveloper/blog/internal/logger"
	"github.com/aveloper/blog/internal/server"
	"github.com/spf13/cobra"
)

func start(cmd *cobra.Command, _ []string) {
	apiOnly, err := cmd.Flags().GetBool("api-only")
	cobra.CheckErr(err)

	cfg := config.Get()
	log := logger.New(&logger.Config{Production: cfg.Production})

	srv := server.NewServer(&server.Config{
		Port:    cfg.Port,
		Logger:  log,
		APIOnly: apiOnly,
	})

	go srv.Listen()

	srv.WaitForShutdown()

	log.Info("Server Shutdown complete")
}
