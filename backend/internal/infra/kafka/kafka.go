package kafka

import (
	"time"

	"github.com/segmentio/kafka-go"
)

// NewReader creates a new Kafka reader configured for the given consumer group,
// topic, and broker list. It starts reading from the earliest available offset
// and commits offsets at a fixed interval.
func NewReader(groupID string, topic string, brokers []string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        groupID,
		Topic:          topic,
		StartOffset:    kafka.FirstOffset,
		CommitInterval: time.Second, // Disable auto-commit
		MaxBytes:       10e6,        // 10 MB
	})
}
