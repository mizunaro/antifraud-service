package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	"github.com/mizunaro/antifraud-service/internal/domain"
	"github.com/mizunaro/antifraud-service/internal/repository"
)

type Worker struct {
	consumer     *Consumer
	postgresRepo *repository.PostgresRepo
	redisRepo    *repository.RedisRepo
}

func NewWorker(c *Consumer, p *repository.PostgresRepo, r *repository.RedisRepo) *Worker {
	return &Worker{consumer: c, postgresRepo: p, redisRepo: r}
}

func (w *Worker) Start(ctx context.Context) error {
	log.Info().Msg("Worker started")

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		msg, err := w.consumer.Fetch(ctx)
		if err != nil {
			return err
		}

		if err := w.processMessage(ctx, msg); err != nil {
			log.Error().Err(err).Msg("failed to process message, skipping...")
			continue
		}
	}
}

func (w *Worker) processMessage(ctx context.Context, msg kafka.Message) error {
	var check domain.URLCheck
	if err := json.Unmarshal(msg.Value, &check); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	status, _ := w.redisRepo.GetStatus(ctx, check.URL)

	if status == domain.StatusPending {
		log.Info().Str("url", check.URL).Msg("Cache miss, analyzing...")

		time.Sleep(1 * time.Second)
		status = domain.URLStatus(rand.IntN(2) + 1)

		if err := w.redisRepo.SetStatus(ctx, check.URL, status, 24*time.Hour); err != nil {
			log.Warn().Err(err).Msg("failed to save status to redis")
		}
	} else {
		log.Info().Str("url", check.URL).Msg("Cache hit! Skipping analysis")
	}

	if err := w.postgresRepo.UpdateStatus(ctx, check.ID, status); err != nil {
		return fmt.Errorf("postgresRepo.UpdateStatus: %w", err)
	}

	log.Info().
		Str("id", check.ID.String()).
		Int("status", int(status)).
		Msg("URL check completed and saved")

	return w.consumer.Commit(ctx, msg)
}
