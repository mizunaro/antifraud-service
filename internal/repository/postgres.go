package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func NewPostgresDB(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse DSN: %w", err)
	}

	var pool *pgxpool.Pool
	// Ретрай-логика: пробуем подключиться 5 раз
	for i := 0; i < 5; i++ {
		pool, err = pgxpool.NewWithConfig(ctx, cfg)
		if err == nil {
			err = pool.Ping(ctx)
			if err == nil {
				break // Успех!
			}
		}
		
		log.Warn().Msgf("Failed to connect to Postgres, retrying... (%d/5)", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("postgres connection failed after retries: %w", err)
	}

	return pool, nil
}
