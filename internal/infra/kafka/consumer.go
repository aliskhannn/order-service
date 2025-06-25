package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(topic string, brokers []string) *Consumer {
	return &Consumer{
		reader: NewReader(topic, brokers),
	}
}

func (c *Consumer) ConsumeMessage(ctx context.Context, handler func(msg []byte) error) error {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			return fmt.Errorf("error reading message from producer: %w", err)
		}

		if err = handler(m.Value); err != nil {
			log.Printf("error handling message: %v", err)
			continue
		}

		log.Printf("message at offset %d handled successfully: %s\n", m.Offset, string(m.Value))
	}
}
