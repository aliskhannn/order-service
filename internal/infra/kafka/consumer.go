package kafka

import (
	"context"
	"errors"
	"github.com/segmentio/kafka-go"
	"log"
	"sync"
)

type messageHandler interface {
	HandleMessage(ctx context.Context, msg []byte) error
}

type Consumer struct {
	reader  *kafka.Reader
	handler messageHandler
}

func NewConsumer(topic string, addr string, handler messageHandler) *Consumer {
	return &Consumer{
		reader:  NewReader(topic, addr),
		handler: handler,
	}
}

func (c *Consumer) ConsumeMessage(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				log.Println("consumer stopped")
				break
			}
			log.Printf("error reading message: %v", err)
			continue
		}

		if err = c.handler.HandleMessage(ctx, m.Value); err != nil {
			log.Printf("error handling message: %v", err)
			continue
		}

		log.Printf("message at offset %d handled successfully: %s\n", m.Offset, string(m.Value))
	}

	log.Println("ConsumeMessage finished")
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
