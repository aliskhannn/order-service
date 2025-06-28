package kafka

import (
	"github.com/segmentio/kafka-go"
	"time"
)

func NewReader(groupID string, topic string, addr string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{addr},
		GroupID:        groupID,
		Topic:          topic,
		MaxBytes:       10e6,        // 10 MB
		CommitInterval: time.Second, // flush commits to kafka every second
	})
}
