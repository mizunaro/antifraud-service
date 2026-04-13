package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mizunaro/antifraud-service/internal/domain"
	kafka "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			Topic:                  topic,
			Balancer:               &kafka.Hash{},
			AllowAutoTopicCreation: true,
		},
	}
}

func (p *Producer) PublishURLCheck(ctx context.Context, check domain.URLCheck) error {
	payload, err := json.Marshal(check)
	if err != nil {
		return fmt.Errorf("marshal check: %w", err)
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(check.ID.String()),
		Value: payload,
	})
	if err != nil {
		return fmt.Errorf("kafka write: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
