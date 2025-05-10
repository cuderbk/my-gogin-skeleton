package kafka

import (
	"context"
	"errors"

	"my-gogin-skeleton/internal/common/logger"

	"github.com/IBM/sarama"
)

type TopicMeta struct {
	Name        string
	Concurrency int
}

type Consumer struct {
	group    sarama.ConsumerGroup
	registry Registry
	topics   []TopicMeta
}

func NewConsumer(
	brokers []string,
	groupID string,
	topics []TopicMeta,
	reg Registry,
	base SaramaBase,
) (*Consumer, error) {

	if len(topics) == 0 {
		logger.Logger.Error("no topics provided")
		return nil, errors.New("no topics provided")
	}
	if reg == nil {
		logger.Logger.Error("no registry provided")
		return nil, errors.New("no registry provided")
	}
	if len(reg) == 0 {
		logger.Logger.Error("no handlers provided")
		return nil, errors.New("no handlers provided")
	}
	for _, topic := range topics {
		if _, ok := reg[topic.Name]; !ok {
			logger.Logger.Errorf("no handler registered for topic: %s", topic.Name)
			return nil, errors.New("no handler for topic " + topic.Name)
		}
	}

	group, err := sarama.NewConsumerGroup(brokers, groupID, base.Build())
	if err != nil {
		logger.Logger.Errorf("failed to create consumer group: %v", err)
		return nil, err
	}

	logger.Logger.Infow("Kafka consumer group created",
		"group", groupID,
		"topics", topics,
	)

	return &Consumer{
		group:    group,
		registry: reg,
		topics:   topics,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	names := make([]string, len(c.topics))
	for i := range c.topics {
		names[i] = c.topics[i].Name
	}

	for {
		if err := c.group.Consume(ctx, names, cgHandler{c}); err != nil {
			logger.Logger.Errorw("Kafka consume error", "err", err)
			return err
		}
		if ctx.Err() != nil {
			logger.Logger.Infow("Kafka consumer context canceled")
			return ctx.Err()
		}
	}
}

func (c *Consumer) Close() error {
	logger.Logger.Infow("closing Kafka consumer")
	return c.group.Close()
}

/* ---- internal handler ----- */

type cgHandler struct{ c *Consumer }

func (h cgHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h cgHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h cgHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	meta := h.meta(claim.Topic())
	if meta == nil {
		logger.Logger.Errorw("missing topic metadata", "topic", claim.Topic())
		return errors.New("topic meta not found")
	}

	logger.Logger.Infow("consuming Kafka topic",
		"topic", claim.Topic(),
		"partition", claim.Partition(),
		"initial_offset", claim.InitialOffset(),
		"concurrency", meta.Concurrency,
	)

	for msg := range claim.Messages() {
		logger.Logger.Debugw("message received",
			"topic", msg.Topic,
			"partition", msg.Partition,
			"offset", msg.Offset,
			"key", string(msg.Key),
		)

		if fn, ok := h.c.registry[msg.Topic]; ok {
			err := fn(sess.Context(), msg.Key, msg.Value)
			if err != nil {
				logger.Logger.Errorw("handler error",
					"topic", msg.Topic,
					"offset", msg.Offset,
					"err", err,
				)
				// TODO: Retry or push to DLQ
			}
		} else {
			logger.Logger.Warnw("no handler found for topic",
				"topic", msg.Topic,
			)
		}

		sess.MarkMessage(msg, "")
	}

	return nil
}

func (h cgHandler) meta(topic string) *TopicMeta {
	for i := range h.c.topics {
		if h.c.topics[i].Name == topic {
			return &h.c.topics[i]
		}
	}
	return nil
}
