package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/mizunaro/antifraud-service/internal/domain"
)

type URLRepository interface {
	Save(ctx context.Context, check domain.URLCheck) error
}

type URLProducer interface {
	PublishURLCheck(ctx context.Context, check domain.URLCheck) error
}

type Service struct {
	repo     URLRepository
	producer URLProducer
}

func New(repo URLRepository, producer URLProducer) *Service {
	return &Service{repo: repo, producer: producer}
}

func (s *Service) ProcessURL(ctx context.Context, rawURL string) (domain.URLCheck, error) {
	check := domain.URLCheck{
		ID:        uuid.New(),
		URL:       rawURL,
		Status:    domain.StatusPending,
		CreatedAt: time.Now(),
	}

	err := s.repo.Save(ctx, check)
	if err != nil {
		return domain.URLCheck{}, fmt.Errorf("repo.Save: %w", err)
	}

	err = s.producer.PublishURLCheck(ctx, check)
	if err != nil {
		return domain.URLCheck{}, fmt.Errorf("producer.PublishURLCheck: %w", err)
	}

	return check, nil
}
