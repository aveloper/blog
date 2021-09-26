package main

import (
	"blog/config"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

//server houses the http.Server and other variables for our HTTP server
type server struct {
	http.Server

	logger *zap.Logger
	router *mux.Router

	//connClose channel is closed when the http.Server is shutdown.
	// It can be used to listen when the server closes
	connClose chan int
}

func NewServer(cfg *config.App, logger *zap.Logger) *server {
	r := mux.NewRouter().StrictSlash(true)
	return &server{
		logger:    logger,
		router:    r,
		connClose: make(chan int, 1),
		Server: http.Server{
			Addr:         fmt.Sprintf("%s:%d", "", cfg.Port),
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
	}
}

func (s *server) Initialize() {
	defer s.graceFullShutdown()

	addRoutes(s.router)
	s.Handler = s.router
}

func (s *server) graceFullShutdown() {
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM)

		sig := <-sigint
		s.logger.Info("OS terminate signal received", zap.String("signal", sig.String()))

		s.logger.Debug("Shutting down server")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := s.Shutdown(ctx)
		if err != nil {
			s.logger.Error("Error shutting down server", zap.Error(err))
		}

		close(s.connClose)
	}()
}

func (s *server) Listen() {
	s.logger.Info("Starting server...", zap.String("address", s.Addr))
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		s.logger.Fatal("HTTP server error", zap.Error(err))
	}
}
