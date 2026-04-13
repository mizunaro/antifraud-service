package app

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/mizunaro/antifraud-service/internal/config"
	"github.com/mizunaro/antifraud-service/internal/repository"
	"github.com/mizunaro/antifraud-service/internal/service"
	http_transport "github.com/mizunaro/antifraud-service/internal/transport/http"
	"github.com/mizunaro/antifraud-service/internal/transport/kafka"
)

func Run(ctx context.Context, c *config.Config) error {
	// Инициализация ресурсов (Postgres, Kafka)
	repo, err := repository.NewPostgresDB(ctx, c.Postgres.DSN)
	if err != nil {
		return err
	}
	defer repo.Close()

	producer := kafka.NewProducer(c.Kafka.Brokers, c.Kafka.Topic)
	defer producer.Close()

	svc := service.New(repo, producer)

	h := http_transport.NewHandler(svc)

	// HTTP Сервер
	router := chi.NewRouter()
	h.Register(router)

	srv := &http.Server{
		Addr:    ":" + c.HTTP.Port,
		Handler: router,
	}

	// Запускаем сервер в горутине
	go func() {
		log.Info().Msgf("server started on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("http server error")
		}
	}()

	// Ждем сигнала отмены из main (через контекст)
	<-ctx.Done()
	log.Info().Msg("shutting down app...")

	// 5. Даем 5 секунд на очистку ресурсов
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("server shutdown failed")
	}

	return nil
}
