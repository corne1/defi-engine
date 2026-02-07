package main

import (
	"context"
	"os"
	"time"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"github.com/corne1/defi-engine/internal/app/dto"
	"github.com/corne1/defi-engine/internal/state"

	"github.com/corne1/defi-engine/internal/app/config"
	"github.com/corne1/defi-engine/internal/observability/logging"
	"github.com/corne1/defi-engine/internal/storage/postgres"
)

func main() {
	// 1️⃣ Context для graceful shutdown
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// 2️⃣ Загружаем конфиг и создаём логгер
	cfg := config.Load()
	logger := logging.New(cfg.Env)

	// 3️⃣ Подключаем Postgres
	db, err := postgres.New(ctx, cfg.DB)
	if err != nil {
		logger.Error("failed to connect to postgres", "err", err)
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("postgres connected")

	// 4️⃣ Инициализация репозитория транзакций
	txRepo := postgres.NewTransactionRepository(db)

	logger.Info("worker started")

	// 5️⃣ Основной цикл воркера
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("worker shutting down")
			return

		case <-ticker.C:
			process(ctx, logger, txRepo)
		}
	}
}

func process(ctx context.Context, logger *slog.Logger, repo *postgres.TransactionRepository) {
	logger.Info("worker tick - processing pending transactions")

	// Здесь пока простая имитация
	pendingTxs := []dto.Transaction{
		{
			ID:    uuid.New(),
			Hash:  "",
			State: state.TxStatePending,
		},
	}

	for _, tx := range pendingTxs {
		// пробуем перевести в sent
		nextState, err := state.Transition(tx.State, state.TxStateSent)
		if err != nil {
			logger.Error("failed state transition", "txID", tx.ID, "err", err)
			continue
		}

		tx.State = nextState

		// сохраняем в БД
		err = repo.UpdateState(ctx, tx.ID, state.TxStatePending, state.TxStateSent)
		if err != nil {
			logger.Error("failed to update state in db", "txID", tx.ID, "err", err)
			continue
		}

		logger.Info("transaction moved to sent", "txID", tx.ID)
	}
}	