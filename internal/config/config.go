package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type (
	Config struct {
		App      AppConfig
		HTTP     HTTPConfig
		Postgres PostgresConfig
		Kafka    KafkaConfig
	}

	AppConfig struct {
		Name    string `envconfig:"APP_NAME" default:"antifraud-service"`
		Version string `envconfig:"APP_VERSION" default:"1.0.0"`
	}

	HTTPConfig struct {
		Port string `envconfig:"HTTP_PORT" default:"8080"`
	}

	PostgresConfig struct {
		DSN string `envconfig:"PG_DSN" required:"true"`
	}

	KafkaConfig struct {
		Brokers []string `envconfig:"KAFKA_BROKERS" default:"localhost:9092"`
		Topic   string   `envconfig:"KAFKA_TOPIC"   default:"urls_to_check"`
	}
)

func New() (*Config, error) {
	var cfg Config

	_ = godotenv.Load(".env")

	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("envconfig: %w", err)
	}
	return &cfg, nil
}
