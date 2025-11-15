package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func ConsumeImageTasks(ctx context.Context, consumer *kafka.Reader, handler func(imagePath string) error) {
	for {
		message, err := consumer.ReadMessage(ctx)
		if err != nil {
			log.Printf("[kafka] error reading message: %v", err)
			continue
		}

		imagePath := string(message.Value)
		err = handler(imagePath)
		if err != nil {
			log.Printf("[kafka]error processing image %s: %v", imagePath, err)
		} else {
			log.Printf("[kafka] iamge successfully processed: %s", imagePath)
		}
	}
}
