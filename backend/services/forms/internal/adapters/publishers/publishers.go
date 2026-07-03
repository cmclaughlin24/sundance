package publishers

import (
	"fmt"
	"log/slog"
	"sundance/backend/services/forms/internal/core/ports"
)

type PublisherType string

const (
	PublisherTypeInMemory PublisherType = "in-memory"
	PublisherTypeKafka    PublisherType = "kafka"
)

type PublisherOptions struct {
	Type  PublisherType `json:"type" env:"TYPE"`
	Kafka *KafkaOptions `json:"kafka,omitempty" envPrefix:"KAFKA_" env:",init"`
}

func Bootstrap(logger *slog.Logger, options PublisherOptions) (ports.Publisher, error) {
	switch options.Type {
	case PublisherTypeInMemory:
		return NewInMemoryPublisher(logger), nil
	case PublisherTypeKafka:
		return NewKafkaPublisher(logger, options.Kafka), nil
	default:
		return nil, fmt.Errorf("unknown publisher type: %s", options.Type)
	}
}
