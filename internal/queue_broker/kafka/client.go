package kafka

import "fmt"

type Client struct {
	*Producer
	*Consumer
}

func NewKafkaClient(cfg *Config) (*Client, error) {
	producer, err := NewProducer(cfg)
	if err != nil {
		return nil, fmt.Errorf("[kafka-client] error init kafka-producer: %w", err)
	}
	consumer, err := NewConsumer(cfg)
	if err != nil {
		producer.Close()
		return nil, fmt.Errorf("[kafka-client] error init kafka-consumer: %w", err)
	}
	return &Client{
		Producer: producer,
		Consumer: consumer,
	}, nil
}

func (c *Client) Close() error {
	if c.Producer != nil {
		c.Producer.Close()
	}
	if c.Consumer != nil {
		return c.Consumer.Close()
	}
	return nil
}
