package config

type KafkaConfig struct {
	Brokers     []string `yaml:"brokers"`
	ClientID    string   `yaml:"client_id"`
	Retries     int      `yaml:"retries"`
	Compression string   `yaml:"compression"` // e.g., snappy, gzip, none
	Acks        string   `yaml:"acks"`        // all, 1, 0
}

func (c *Config) KafkaConfig() KafkaConfig {
	return c.Kafka
}
