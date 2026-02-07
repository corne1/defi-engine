package main

import (
	"path/filepath"

	"github.com/corne1/defi-engine/internal/app/config"
	"github.com/corne1/defi-engine/internal/observability/logging"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// 1️⃣ Загружаем конфиг и создаём логгер
	cfg := config.Load()
	logger := logging.New(cfg.Env)

	logger.Info("starting migrations...")

	// 2️⃣ DSN к базе
	dsn := cfg.DB.DSN()
	absPath, err := filepath.Abs("./migrations")
	if err != nil {
		logger.Error("failed to get absolute path", "err", err)
		return
	}
	// 3️⃣ Инициализация мигратора
	absPath = filepath.ToSlash(absPath)

	// ⚡ Windows-safe
	m, err := migrate.New(
		"file://./migrations",
		dsn,
	)

	if err != nil {
		logger.Error("failed to init migrate", "err", err)
		return
	}

	// 4️⃣ Применяем миграции
	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			logger.Info("no new migrations to apply")
		} else {
			logger.Error("migration failed", "err", err)
			return
		}
	} else {
		logger.Info("migrations applied successfully")
	}
}
