package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"

	"github.com/mizunaro/antifraud-service/internal/config"
	"github.com/mizunaro/antifraud-service/internal/repository"
	"github.com/mizunaro/antifraud-service/internal/service"
	kafka_transport "github.com/mizunaro/antifraud-service/internal/transport/kafka"
)

func RunWorker(ctx context.Context, c *config.Config) error {
	pgRepo, err := repository.NewPostgresDB(ctx, c.Postgres.DSN)
	if err != nil {
		return err
	}
	defer pgRepo.Close()

	consumer := kafka_transport.NewConsumer(c.Kafka.Brokers, c.Kafka.Topic, c.Kafka.GroupID)
	defer consumer.Close()

	redisRepo := repository.NewRedis(c.Redis.Addr, c.Redis.Password, c.Redis.DB)
	defer redisRepo.Close()

	w := service.NewWorker(consumer, pgRepo, redisRepo, c.Worker.BadWords)

	log.Info().Msg("worker: starting to consume messages...")

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Info().Msg("Metrics server started on :9090")
		if err := http.ListenAndServe(":9090", nil); err != nil {
			log.Error().Err(err).Msg("metrics server error")
		}
	}()

	if err := w.Start(ctx); err != nil {
		return fmt.Errorf("worker run error: %w", err)
	}

	log.Info().Msg("worker: stopped gracefully")
	return nil
}
