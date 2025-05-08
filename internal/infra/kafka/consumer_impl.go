package kafka

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

func (c *kafkaConsumer) Start(ctx context.Context) error {
	go func() {
		for {
			if err := c.client.Consume(ctx, c.topics, c); err != nil {
				log.Printf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()
	return nil
}

func (c *kafkaConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	sem := make(chan struct{}, c.concurrency) // worker pool

	for msg := range claim.Messages() {
		sem <- struct{}{}
		go func(m *sarama.ConsumerMessage) {
			defer func() { <-sem }()
			ctx := context.WithValue(context.Background(), "trace_id", string(m.Key))
			err := c.handler.HandleMessage(ctx, m.Topic, string(m.Key), m.Value)
			if err == nil {
				sess.MarkMessage(m, "")
			} else {
				log.Printf("error handling message: %v", err)
			}
		}(msg)
	}
	return nil
}
