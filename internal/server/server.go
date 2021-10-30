package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aveloper/blog/internal/server/api"
	"github.com/aveloper/blog/internal/server/web"

	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

//Server houses the http.Server and other variables for our HTTP server
type Server struct {
	server http.Server

	logger *zap.Logger
	router *mux.Router

	//connClose channel is closed when the http.Server is shutdown.
	// It can be used to listen when the server closes
	connClose chan int

	apiOnly bool
}

type Config struct {
	Port    int
	Logger  *zap.Logger
	APIOnly bool
}

//NewServer creates a new instance of Server
func NewServer(cfg *Config) *Server {
	r := mux.NewRouter().StrictSlash(true)
	return &Server{
		logger:    cfg.Logger,
		router:    r,
		connClose: make(chan int, 1),
		server: http.Server{
			Addr:         fmt.Sprintf("%s:%d", "", cfg.Port),
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
		apiOnly: cfg.APIOnly,
	}
}

func (s *Server) Listen() {
	s.setup()

	s.logger.Info("Starting server...", zap.String("address", s.server.Addr))
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		s.logger.Fatal("HTTP server error", zap.Error(err))
	}
}

func (s *Server) WaitForShutdown() {
	<-s.connClose
}

func (s *Server) setup() {
	defer s.graceFullShutdown()

	if !s.apiOnly {
		web.Routes(s.router)
	}

	apiRouter := s.router.PathPrefix("/api").Subrouter()
	api.Routes(apiRouter)

	// Add middlewares and handlers here

	s.server.Handler = s.router
}

func (s *Server) graceFullShutdown() {
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM)

		sig := <-sigint
		s.logger.Info("OS terminate signal received", zap.String("signal", sig.String()))

		s.logger.Debug("Shutting down server")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := s.server.Shutdown(ctx)
		if err != nil {
			s.logger.Error("Error shutting down server", zap.Error(err))
		}

		close(s.connClose)
	}()
}