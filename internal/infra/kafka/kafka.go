package kafka

import "github.com/segmentio/kafka-go"

func NewReader(topic string, brokers []string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		MaxBytes: 10e6, // 10 MB
	})
}

func NewWriter(topic string, brokers []string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}
