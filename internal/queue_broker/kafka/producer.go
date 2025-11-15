package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func ProduceImageTask(ctx context.Context, producer *kafka.Writer, imagePath string) error {
	message := kafka.Message{
		Key:   []byte(imagePath),
		Value: []byte(imagePath),
	}
	err := producer.WriteMessages(ctx, message)
	if err != nil {
		log.Printf("[kafka] failed to produce message: %v", err)
		return err
	}
	return nil
}
