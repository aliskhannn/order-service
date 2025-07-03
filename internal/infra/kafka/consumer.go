package kafka

import (
	"context"
	"errors"
	"sync"

	"go.uber.org/zap"

	"github.com/segmentio/kafka-go"

	customerr "github.com/aliskhannn/order-service/internal/errors"
)

type messageHandler interface {
	HandleMessage(ctx context.Context, msg []byte) error
}

type Consumer struct {
	reader  *kafka.Reader
	logger  *zap.Logger
	handler messageHandler
}

func NewConsumer(groupID string, topic string, brokers []string, l *zap.Logger, h messageHandler) *Consumer {
	return &Consumer{
		reader:  NewReader(groupID, topic, brokers),
		logger:  l,
		handler: h,
	}
}

func (c *Consumer) ConsumeMessage(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	defer c.Close()

	go func() {
		<-ctx.Done()
		c.logger.Info("Received shutdown signal, closing consumer")
		c.Close()
	}()

	c.logger.Info("ConsumeMessage started",
		zap.String("topic", c.reader.Config().Topic),
		zap.String("groupID", c.reader.Config().GroupID),
	)

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				c.logger.Warn("context canceled or deadline exceeded, stopping consumer", zap.Error(err))
				break
			}

			c.logger.Error("error reading message", zap.Error(err))
			continue
		}

		if err = c.handler.HandleMessage(ctx, m.Value); err != nil {
			c.handleMessageError(m, err)
			continue
		}

		c.logger.Info("message handled successfully",
			zap.Int64("offset", m.Offset),
			zap.String("message", string(m.Value)),
		)
	}

	c.logger.Info("ConsumeMessage finished")
}

func (c *Consumer) handleMessageError(m kafka.Message, err error) {
	msgStr := string(m.Value)

	switch {
	case errors.Is(err, customerr.ErrInvalidJSON):
		c.logger.Warn("invalid JSON format",
			zap.Int64("offset", m.Offset),
			zap.String("message", msgStr),
			zap.Error(err),
		)
	case errors.Is(err, customerr.ErrNilOrder):
		c.logger.Warn("nil order received",
			zap.Int64("offset", m.Offset),
			zap.String("message", msgStr),
			zap.Error(err),
		)
	case errors.Is(err, customerr.ErrValidation):
		c.logger.Warn("validation error",
			zap.Int64("offset", m.Offset),
			zap.String("message", msgStr),
			zap.Error(err),
		)
	case errors.Is(err, customerr.ErrCreateOrder):
		c.logger.Warn("failed to create order",
			zap.Int64("offset", m.Offset),
			zap.String("message", msgStr),
			zap.Error(err),
		)
	default:
		c.logger.Error("unexpected error while handling message",
			zap.Int64("offset", m.Offset),
			zap.String("message", msgStr),
			zap.Error(err),
		)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
