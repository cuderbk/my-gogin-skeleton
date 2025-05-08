package kafka

import (
	"context"
	"fmt"
)

type TopicHandlerMap map[string]ConsumerHandler

func (h TopicHandlerMap) HandleMessage(ctx context.Context, topic string, key string, value []byte) error {
	if handler, ok := h[topic]; ok {
		return handler.HandleMessage(ctx, topic, key, value)
	}
	return fmt.Errorf("no handler for topic %s", topic)
}
