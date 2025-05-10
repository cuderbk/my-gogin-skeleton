package kafka

import (
	"strings"
	"time"

	"github.com/IBM/sarama"
)

type SaramaBase struct {
	ClientID      string
	Retries       int
	Compression   string
	Acks          string
	CommitTimeout int // ms
}

func (b SaramaBase) Build() *sarama.Config {
	c := sarama.NewConfig()
	c.Version = sarama.V2_8_0_0
	c.ClientID = b.ClientID
	c.Producer.Retry.Max = b.Retries
	c.Producer.RequiredAcks = toAcks(b.Acks)
	c.Producer.Compression = toCodec(b.Compression)
	c.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	c.Consumer.Offsets.AutoCommit.Interval = time.Duration(b.CommitTimeout) * time.Millisecond
	return c
}

func toAcks(in string) sarama.RequiredAcks {
	switch strings.ToLower(in) {
	case "0", "none":
		return sarama.NoResponse
	case "1", "leader":
		return sarama.WaitForLocal
	default:
		return sarama.WaitForAll
	}
}
func toCodec(in string) sarama.CompressionCodec {
	switch strings.ToLower(in) {
	case "gzip":
		return sarama.CompressionGZIP
	case "lz4":
		return sarama.CompressionLZ4
	case "zstd":
		return sarama.CompressionZSTD
	case "snappy":
		return sarama.CompressionSnappy
	default:
		return sarama.CompressionNone
	}
}
