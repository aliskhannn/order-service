package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aliskhannn/order-service/internal/model"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(topic string, brokers []string) *Producer {
	return &Producer{
		writer: NewWriter(topic, brokers),
	}
}

func (p *Producer) ProduceMessage(ctx context.Context, order model.Order) error {
	b, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("error narshaling data: %w", err)
	}

	return p.writer.WriteMessages(ctx, kafka.Message{Value: b})
}
