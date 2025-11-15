package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

type Consumer struct {
	group sarama.ConsumerGroup
	topic string
}

func NewConsumer(cfg *Config) (*Consumer, error) {
	group, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, cfg.SaramaConfig())
	if err != nil {
		return nil, err
	}

	return &Consumer{
		group: group,
		topic: cfg.Topic,
	}, nil
}

// consumerGroupHandler оборачивает функцию обработки сообщений
type consumerGroupHandler struct {
	handler func(string) error
}

func (h *consumerGroupHandler) Setup(s sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerGroupHandler) Cleanup(s sarama.ConsumerGroupSession) error { return nil }
func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if err := h.handler(string(msg.Value)); err != nil {
			log.Printf("[kafka-consumer] handler error: %v", err)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

func (c *Consumer) Consume(ctx context.Context, handler func(string) error) error {
	h := &consumerGroupHandler{handler: handler}
	for {
		err := c.group.Consume(ctx, []string{c.topic}, h)
		if err != nil {
			log.Printf("[kafka-consumer] error during consuming: %v", err)
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func (c *Consumer) Close() error {
	if c.group != nil {
		return c.group.Close()
	}
	return nil
}
