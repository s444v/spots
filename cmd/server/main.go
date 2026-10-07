package main

import (
	"context"
	"errors"
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
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Error("CFG error", "error", err)
		os.Exit(1)
	}

	s := server.New(cfg, webFS, log)

	log.Info("server starting")

	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Info("server starting error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := s.Shutdown(shutdownCtx); err != nil {
		log.Info("forced shutdown", "error", err)
	}
	log.Info("server shutdown")

}
