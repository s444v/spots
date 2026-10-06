package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/s444v/spots/internal/config"
	"github.com/s444v/spots/internal/logger"
	"github.com/s444v/spots/internal/server"
	"github.com/s444v/spots/web"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.Init()

	webFS, err := web.Static()
	if err != nil {
		logger.Log.Error("server stopped", "app", "spots", "error", err)
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		logger.Log.Error("server stopped", "app", "spots", "error", err)
		os.Exit(1)
	}

	s := server.NewServer(cfg, webFS)
	if err != nil {
		logger.Log.Error("server stopped", "app", "spots", "error", err)
		os.Exit(1)
	}

	logger.Log.Info("server starting", "app", "spots")

	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Info("server error", "app", "spots", "error", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Shutdown(shutdownCtx); err != nil {
		logger.Log.Info("forced shutdown", "app", "spots", "error", err)
	}
	logger.Log.Info("server shutdown", "app", "spots")

}
