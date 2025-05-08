package main

import (
	"context"
	"fmt"

	"logging/config"
	"logging/internal/infra/cache"
	"logging/internal/infra/db"
)

type Infra struct {
	DB    db.Store
	Redis cache.Store
	Close func()
}

func InitInfra(ctx context.Context, cfg *config.Config) (*Infra, error) {
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

	// Compose closer
	cleanup := func() {
		dbConn.Close()
		redisClient.Close()
	}

	return &Infra{
		DB:    db.NewPostgresStore(dbConn),
		Redis: cache.NewRedisStore(redisClient),
		Close: cleanup,
	}, nil
}
