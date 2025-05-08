package kafka

import "github.com/IBM/sarama"

type kafkaConsumer struct {
	client      sarama.ConsumerGroup
	groupID     string
	topics      []string
	handler     ConsumerHandler
	concurrency int
}

func NewConsumer(brokers []string, groupID string, topics []string, handler ConsumerHandler, concurrency int) (Consumer, error) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_6_0_0
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	client, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
	if err != nil {
		return nil, err
	}

	return &kafkaConsumer{
		client:      client,
		groupID:     groupID,
		topics:      topics,
		handler:     handler,
		concurrency: concurrency,
	}, nil
}
