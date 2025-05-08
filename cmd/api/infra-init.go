package main

import (
	"context"
	"fmt"

	"logging/config"
	"logging/internal/infra/cache"
	"logging/internal/infra/db"
	"logging/internal/infra/kafka"
)

type Infra struct {
	DB       db.Store
	Redis    cache.Store
	Consumer kafka.Consumer
	errCh    chan error
	Close    func()
}

func InitInfra(ctx context.Context, cfg *config.Config, handlers kafka.TopicHandlerMap) (*Infra, error) {
	// Connect Postgres
	dbConn, err := db.InitDB(ctx, cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Postgres: %w", err)
	}

	// Connect Redis
	redisClient, err := cache.InitRedis(cfg.Redis)
	if err != nil {
		dbConn.Close()
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	consumer, err := kafka.NewConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.GroupID,
		handlers.Topics(),
		topicHandlers,
		8, // concurrency level
	)
	if err != nil {
		dbConn.Close()
		redisClient.Close()
		return nil, fmt.Errorf("failed to init kafka consumer: %w", err)
	}

	errCh := make(chan error, 1)

	// Compose closer
	cleanup := func() {
		dbConn.Close()
		redisClient.Close()
	}

	return &Infra{
		DB:       db.NewPostgresStore(dbConn),
		Redis:    cache.NewRedisStore(redisClient),
		Consumer: consumer,
		errCh:    errCh,
		Close:    cleanup,
	}, nil
}

func (i *Infra) Start(ctx context.Context) {
	go func() {
		if err := i.consumer.Start(ctx); err != nil {
			i.errCh <- err
		}
	}()
}

func (i *Infra) Wait() error {
	select {
	case <-ctx.Done():
		return nil
	case err := <-i.errCh:
		return err
	}
}

func (i *Infra) Stop() {
	i.consumer.Stop()
	i.Close()
}
