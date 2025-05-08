package kafka

import "context"

type ConsumerHandler interface {
	HandleMessage(ctx context.Context, topic string, key string, value []byte) error
}
