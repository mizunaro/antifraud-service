package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mizunaro/antifraud-service/internal/config"
	"github.com/mizunaro/antifraud-service/internal/worker"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config init failed: %v", err)
	}

	// Настраиваем логгер
	// logger.Init(cfg.Logger)

	// Создаем контекст для Graceful Shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Запускаем основной цикл приложения
	if err := worker.Run(ctx, cfg); err != nil {
		log.Fatalf("worker finished with error: %v", err)
	}
}
