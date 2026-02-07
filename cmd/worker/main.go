package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"log/slog"

	"github.com/corne1/defi-engine/internal/app/config"
	"github.com/corne1/defi-engine/internal/observability/logging"
	"github.com/corne1/defi-engine/internal/storage/postgres"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	cfg := config.Load()
	logger := logging.New(cfg.Env)

	db, err := postgres.New(ctx, cfg.DB)
	if err != nil {
		logger.Error("failed to connect to postgres", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("worker started")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("worker shutting down")
			return

		case <-ticker.C:
			process(ctx, logger)
		}
	}
}

func process(ctx context.Context, logger *slog.Logger) {
	logger.Info("worker tick")
	// здесь будет:
	// - поиск pending tx
	// - отправка в EVM
	// - обновление state
}