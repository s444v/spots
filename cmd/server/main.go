package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/s444v/spots/internal/config"
	"github.com/s444v/spots/internal/logger"
	"github.com/s444v/spots/internal/server"
	"github.com/s444v/spots/web"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log := logger.Init()

	webFS, err := web.Static()
	if err != nil {
		log.Error("FS error", "error", err)
		return
	}

	cfg, err := config.Load()
	if err != nil {
		log.Error("CFG error", "error", err)
		return
	}

	s := server.New(cfg, webFS, log)

	ln, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		log.Error("listen failed", "addr", cfg.Addr, "err", err)
		return
	}
	log.Info("server listening", "addr", ln.Addr().String())

	errChan := make(chan error, 1)

	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Info("shutting down server gracefully...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		if err := s.Shutdown(shutdownCtx); err != nil {
			log.Error("forced shutdown", "error", err)
		}
		log.Info("server shutdown complete")

	case err := <-errChan:
		log.Error("server starting error", "error", err)
	}

}
