package config

type TopicConfig struct {
	Name        string `mapstructure:"name" yaml:"name"`
	Concurrency int    `mapstructure:"concurrency" yaml:"concurrency"`
}

type KafkaConfig struct {
	Brokers       []string      `mapstructure:"brokers" yaml:"brokers"`
	ClientID      string        `mapstructure:"client_id" yaml:"client_id"`
	Retries       int           `mapstructure:"retries" yaml:"retries"`
	Compression   string        `mapstructure:"compression" yaml:"compression"` // snappy, gzip, none
	Acks          string        `mapstructure:"acks" yaml:"acks"`               // all, 1, 0
	ConsumerGroup string        `mapstructure:"consumer_group" yaml:"consumer_group"`
	CommitTimeout int           `mapstructure:"commit_timeout" yaml:"commit_timeout"` // milliseconds
	Topics        []TopicConfig `mapstructure:"topics" yaml:"topics"`
}

/*  Getter  */
func (c *Config) KafkaCfg() KafkaConfig { return c.Kafka }
