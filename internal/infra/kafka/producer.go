package kafka

import "github.com/IBM/sarama"

type Producer struct{ sarama.SyncProducer }

/* Log and DLQ Handling in the future */
// TODO - Add log and DLQ handling
// TODO - Add async producer

func NewProducer(brokers []string, base SaramaBase) (*Producer, error) {
	p, err := sarama.NewSyncProducer(brokers, base.Build())
	if err != nil {
		return nil, err
	}
	return &Producer{p}, nil
}
