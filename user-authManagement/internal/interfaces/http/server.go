package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// Server wraps http.Server with graceful shutdown semantics.
type Server struct {
	srv      *http.Server
	logger   *slog.Logger
	shutdown time.Duration
}

func NewServer(handler http.Handler, port int, readTO, writeTO, shutdownTO time.Duration, logger *slog.Logger) *Server {
	return &Server{
		srv: &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      handler,
			ReadTimeout:  readTO,
			WriteTimeout: writeTO,
			IdleTimeout:  60 * time.Second,
		},
		logger:   logger,
		shutdown: shutdownTO,
	}
}

// Start blocks until ListenAndServe returns. Returns nil on graceful shutdown.
func (s *Server) Start() error {
	s.logger.Info("http server starting", slog.String("addr", s.srv.Addr))
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Stop initiates a graceful shutdown bounded by the configured timeout.
func (s *Server) Stop(ctx context.Context) error {
	sctx, cancel := context.WithTimeout(ctx, s.shutdown)
	defer cancel()
	return s.srv.Shutdown(sctx)
}
