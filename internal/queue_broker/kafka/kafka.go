package kafka

import (
	"log"

	"github.com/segmentio/kafka-go"
)

type Config struct {
	Brokers []string
	Topic   string
	GroupID string
}

// NewProducer созадет kafka producer
func NewProducer(cfg Config) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  cfg.Brokers,
		Topic:    cfg.Topic,
		Balancer: &kafka.LeastBytes{},
	})
}

// NewConsumer создает kafka consumer
func NewConsumer(cfg Config) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		GroupID: cfg.GroupID,
		Topic:   cfg.Topic,
	})
}

// CloseProducer закрывает kafka producer
func CloseProducer(producer *kafka.Writer) {
	err := producer.Close()
	if err != nil {
		log.Printf("[kafka] failed to close producer: %v", err)
	}
}

// CloseConsumer закрывает kafka consumer
func CloseConsumer(consumer *kafka.Reader) {
	err := consumer.Close()
	if err != nil {
		log.Printf("[kafka] failrd tp clode consumer: %v", err)
	}
}
