package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"strconv"
	"time"
	"github.com/corne1/defi-engine/internal/app/config"
	"github.com/corne1/defi-engine/internal/observability/logging"
)

func main() {
	cfg := config.Load()
	logger := logging.New(cfg.Env)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	server := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.API.Port),
		Handler: mux,
	}

	go func() {
		logger.Info("api started", "port", cfg.API.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown failed", "err", err)
	}

	logger.Info("server exited properly")
}
