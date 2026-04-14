package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"

	"github.com/mizunaro/antifraud-service/internal/domain"
	"github.com/mizunaro/antifraud-service/internal/repository"
	kafka_transport "github.com/mizunaro/antifraud-service/internal/transport/kafka"
)

type Worker struct {
	consumer     *kafka_transport.Consumer
	postgresRepo *repository.PostgresRepo
	redisRepo    *repository.RedisRepo
	badWords     []string
}

func NewWorker(
	c *kafka_transport.Consumer,
	p *repository.PostgresRepo,
	r *repository.RedisRepo,
	b []string,
) *Worker {
	return &Worker{consumer: c, postgresRepo: p, redisRepo: r, badWords: b}
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

	isCacheHit := false
	status, _ := w.redisRepo.GetStatus(ctx, check.URL)

	if status == domain.StatusPending {
		status = w.analyze(check.URL)
		_ = w.redisRepo.SetStatus(ctx, check.URL, status, 24*time.Hour)
	} else {
		isCacheHit = true
	}

	processedURLs.WithLabelValues(fmt.Sprintf("%d", status), fmt.Sprintf("%v", isCacheHit)).Inc()

	if err := w.postgresRepo.UpdateStatus(ctx, check.ID, status); err != nil {
		return fmt.Errorf("postgresRepo.UpdateStatus: %w", err)
	}

	log.Info().
		Str("component", "worker").
		Str("id", check.ID.String()).
		Str("url", check.URL).
		Int("status", int(status)).
		Bool("cache_hit", isCacheHit).
		Msg("URL processing finished")

	return w.consumer.Commit(ctx, msg)
}

func (w *Worker) analyze(rawURL string) domain.URLStatus {
	url := strings.ToLower(rawURL)

	for _, badWord := range w.badWords {
		if strings.Contains(url, badWord) {
			return domain.StatusMalicious
		}
	}

	return domain.StatusSafe
}
