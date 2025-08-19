package kafka

import (
	"context"
	"errors"
	"github.com/aliskhannn/order-service/internal/kafka/handlers/order"
	"sync"

	"go.uber.org/zap"

	"github.com/segmentio/kafka-go"
)

type messageHandler interface {
	ProcessMessage(ctx context.Context, msg kafka.Message) error
}

// Consumer wraps a Kafka reader and delegates consumed messages
// to a messageHandler. It also handles logging and graceful shutdown.
type Consumer struct {
	reader  *kafka.Reader
	logger  *zap.Logger
	handler messageHandler
}

// NewConsumer creates a new Kafka consumer bound to a topic, consumer group,
// and a set of broker addresses. A logger is used for structured logs, and
// a messageHandler must be provided for message processing.
func NewConsumer(groupID string, topic string, brokers []string, l *zap.Logger, h messageHandler) *Consumer {
	return &Consumer{
		reader:  NewReader(groupID, topic, brokers),
		logger:  l,
		handler: h,
	}
}

// ConsumeMessage starts consuming messages from Kafka and passes them to
// the configured messageHandler. It runs until the provided context is canceled.
// The WaitGroup is decremented when consumption stops, allowing coordinated shutdown.
//
// This method ensures graceful shutdown by checking for context cancellation.
func (c *Consumer) ConsumeMessage(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		_ = c.Close()
		c.logger.Info("Consumer closed successfully")
	}()

	c.logger.Info("ConsumeMessage started",
		zap.String("topic", c.reader.Config().Topic),
		zap.String("groupID", c.reader.Config().GroupID),
	)

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("shutdown signal received, stopping consumer")
			return
		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					c.logger.Warn("context canceled or deadline exceeded, stopping consumer", zap.Error(err))
					break
				}

				c.logger.Error("error reading message", zap.Error(err))
				continue
			}

			err = c.handler.ProcessMessage(ctx, msg)
			if err != nil {
				if errors.Is(err, order.ErrNilOrder) {
					c.logger.Error("nil order received",
						zap.Int64("offset", msg.Offset),
						zap.String("message", string(msg.Value)),
					)
					continue
				}

				c.logger.Error("error processing message",
					zap.Int64("offset", msg.Offset),
					zap.String("message", string(msg.Value)),
					zap.Error(err),
				)
				continue
			}

			c.logger.Info("message handled successfully",
				zap.Int64("offset", msg.Offset),
				zap.String("message", string(msg.Value)),
			)
		}
	}
}

// Close shuts down the Kafka reader and releases resources.
func (c *Consumer) Close() error {
	return c.reader.Close()
}
