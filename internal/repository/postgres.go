package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/mizunaro/antifraud-service/internal/domain"
)

type PostgresRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresDB(ctx context.Context, dsn string) (*PostgresRepo, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse DSN: %w", err)
	}

	var pool *pgxpool.Pool
	for i := range 5 {
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

	return &PostgresRepo{pool: pool}, nil
}

func (r *PostgresRepo) Save(ctx context.Context, check domain.URLCheck) error {
	query := `
		INSERT INTO url_checks 
			(id, url, status, created_at)
		VALUES
			($1, $2, $3, $4)
		ON CONFLICT (url) DO NOTHING`

	_, err := r.pool.Exec(ctx, query, check.ID, check.URL, check.Status, check.CreatedAt)
	if err != nil {
		return fmt.Errorf("postgres exec: %w", err)
	}

	return nil
}

func (r *PostgresRepo) UpdateStatus(
	ctx context.Context,
	id uuid.UUID,
	status domain.URLStatus,
) error {
	query := `UPDATE url_checks SET status = $1 WHERE id = $2`

	_, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("postgres exec: %w", err)
	}

	return nil
}

func (r *PostgresRepo) Close() {
	r.pool.Close()
}
