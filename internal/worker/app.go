package worker

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/mizunaro/antifraud-service/internal/config"
	"github.com/mizunaro/antifraud-service/internal/repository"
	"github.com/mizunaro/antifraud-service/internal/transport/kafka"
)

func Run(ctx context.Context, c *config.Config) error {
	pgRepo, err := repository.NewPostgresDB(ctx, c.Postgres.DSN)
	if err != nil {
		return err
	}
	defer pgRepo.Close()

	consumer := kafka.NewConsumer(c.Kafka.Brokers, c.Kafka.Topic, c.Kafka.GroupID)
	defer consumer.Close()

	redisRepo := repository.NewRedis(c.Redis.Addr, c.Redis.Password, c.Redis.DB)
	defer redisRepo.Close()

	w := kafka.NewWorker(consumer, pgRepo, redisRepo)

	log.Info().Msg("worker: starting to consume messages...")

	if err := w.Start(ctx); err != nil {
		return fmt.Errorf("worker run error: %w", err)
	}

	log.Info().Msg("worker: stopped gracefully")
	return nil
}
