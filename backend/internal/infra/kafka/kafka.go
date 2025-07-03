package kafka

import (
	"github.com/segmentio/kafka-go"
	"time"
)

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
