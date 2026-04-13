package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/mizunaro/antifraud-service/internal/domain"
)

type RedisRepo struct {
	client *redis.Client
}

func NewRedis(addr string, password string, db int) *RedisRepo {
	return &RedisRepo{
		client: redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: db}),
	}
}

func (r *RedisRepo) SetStatus(
	ctx context.Context,
	url string,
	status domain.URLStatus,
	ttl time.Duration,
) error {
	return r.client.Set(ctx, url, int(status), ttl).Err()
}

func (r *RedisRepo) GetStatus(ctx context.Context, url string) (domain.URLStatus, error) {
	val, err := r.client.Get(ctx, url).Int()
	if err == redis.Nil {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return domain.URLStatus(val), nil
}

func (r *RedisRepo) Close() {
	r.client.Close()
}
