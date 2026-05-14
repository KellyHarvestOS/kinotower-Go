package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/KellyHarvestOS/kinotower-Go/internal/config"
)

type Server struct {
	httpServer *http.Server
	log        *slog.Logger
	timeout    time.Duration
}

func New(cfg config.Config, handler http.Handler, log *slog.Logger) *Server {
	return &Server{
		httpServer: &http.Server{Addr: cfg.AppAddr, Handler: handler, ReadHeaderTimeout: 5 * time.Second},
		log:        log,
		timeout:    cfg.ShutdownTimeout,
	}
}

func (s *Server) Start() error {
	s.log.Info("server listening", "addr", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	return s.httpServer.Shutdown(shutdownCtx)
}
