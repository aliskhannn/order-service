package kafka

import "github.com/segmentio/kafka-go"

func NewReader(topic string, addr string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{addr},
		Topic:    topic,
		MaxBytes: 10e6, // 10 MB
	})
}
