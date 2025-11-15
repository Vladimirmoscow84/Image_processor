package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

type Producer struct {
	topic    string
	producer sarama.SyncProducer
}

func NewProducer(cfg *Config) (*Producer, error) {
	prod, err := sarama.NewSyncProducer(cfg.Brokers, cfg.SaramaConfig())
	if err != nil {
		return nil, err
	}

	return &Producer{
		topic:    cfg.Topic,
		producer: prod,
	}, nil
}

func (p *Producer) Produce(ctx context.Context, msg string) error {
	message := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(msg),
	}
	_, _, err := p.producer.SendMessage(message)
	if err != nil {
		log.Printf("[kafka-producer] failed to send message: %v", err)
		return err
	}
	return nil
}

func (p *Producer) Close() error {
	if p.producer != nil {
		return p.producer.Close()
	}
	return nil
}
