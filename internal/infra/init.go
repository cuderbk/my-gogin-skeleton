package infra

import (
	"context"
	"fmt"

	"logging/config"
	"logging/internal/common/logger"
	"logging/internal/infra/cache"
	"logging/internal/infra/db"
	"logging/internal/infra/kafka"
)

type Infra struct {
	DB     db.Store
	Redis  cache.Store
	Kafka  *kafka.Consumer
	Cancel context.CancelFunc
	Close  func()
}

func InitInfra(ctx context.Context, cfg *config.Config) (*Infra, error) {
	// Connect Postgres
	dbConn, err := db.InitDB(ctx, cfg.DB)
	if err != nil {

		return nil, fmt.Errorf("failed to connect to Postgres: %w", err)
	}

	logger.Logger.Infow("Connected to Postgres")
	dbStore := db.NewPostgresStore(dbConn)
	// Connect Redis
	redisClient, err := cache.InitRedis(cfg.Redis)
	if err != nil {
		dbConn.Close()
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	logger.Logger.Infow("Connected to Redis")
	redisStore := cache.NewRedisStore(redisClient)
	// Connect Kafka
	topicMetas := make([]kafka.TopicMeta, len(cfg.Kafka.Topics))
	for i, t := range cfg.Kafka.Topics {
		topicMetas[i] = kafka.TopicMeta{
			Name:        t.Name,
			Concurrency: t.Concurrency,
		}
	}

	kafkaBase := kafka.SaramaBase{
		ClientID:      cfg.Kafka.ClientID,
		Retries:       cfg.Kafka.Retries,
		Compression:   cfg.Kafka.Compression,
		Acks:          cfg.Kafka.Acks,
		CommitTimeout: cfg.Kafka.CommitTimeout,
	}

	kafkaConsumer, err := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.ConsumerGroup, topicMetas, kafka.DefaultRegistry, kafkaBase)
	if err != nil {
		dbConn.Close()
		redisClient.Close()
		return nil, fmt.Errorf("failed to initialize Kafka consumer: %w", err)
	}
	logger.Logger.Infow("Kafka consumer initialized", "topics", cfg.Kafka.Topics)

	// === Start Kafka Consumer ===
	kafkaCtx, cancel := context.WithCancel(ctx)
	go func() {
		if err := kafkaConsumer.Start(kafkaCtx); err != nil {
			logger.Logger.Infow("Kafka consumer stopped", "err", err)
		}
	}()

	// === Compose Infra ===
	infra := &Infra{
		DB:     dbStore,
		Redis:  redisStore,
		Kafka:  kafkaConsumer,
		Cancel: cancel,
	}

	infra.Close = func() {
		logger.Logger.Infow("Closing Kafka, Redis, Postgres")
		cancel()
		kafkaConsumer.Close()
		redisClient.Close()
		dbConn.Close()
	}

	return infra, nil
}
