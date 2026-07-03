package publishers

import (
	"context"
	"log/slog"
	"strings"
	"sundance/backend/services/forms/internal/core/domain"
	"sundance/backend/services/forms/internal/core/ports"

	"github.com/segmentio/kafka-go"
)

type KafkaOptions struct {
	Addr []string `json:"addr"`
}

type KafkaPublisher struct {
	logger *slog.Logger
	writer *kafka.Writer
}

func NewKafkaPublisher(logger *slog.Logger, options *KafkaOptions) ports.Publisher {
	return &KafkaPublisher{
		logger: logger,
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(options.Addr...),
			AllowAutoTopicCreation: true,
		},
	}
}

func (p *KafkaPublisher) Publish(ctx context.Context, event domain.Event) error {
	message := kafka.Message{
		Topic: topic(event),
		Key:   []byte(event.AggregateID),
		Value: event.Payload,
	}

	if err := p.writer.WriteMessages(ctx, message); err != nil {
		p.logger.ErrorContext(ctx, "failed to publish event", "event_id", event.ID, "error", err)
		return err
	}

	return nil
}

func topic(event domain.Event) string {
	return strings.Join([]string{string(event.AggregateType), string(event.Type)}, ".")
}
