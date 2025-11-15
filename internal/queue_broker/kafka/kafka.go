package kafka

import (
	"context"
)

type producerConsumer interface {
	Produce(ctx context.Context, msg string) error
	Consume(ctx context.Context, handler func(string) error) error
	Close() error
}
