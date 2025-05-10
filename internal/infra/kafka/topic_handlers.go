package kafka

import (
	"context"
	"my-gogin-skeleton/internal/common/logger"
)

func DashboardHandler(_ context.Context, k, v []byte) error { /* todo */

	logger.Logger.Infow("DashboardHandler: key=%s, value=%s", k, v)
	return nil
}
func AlertHandler(_ context.Context, k, v []byte) error { /* todo */

	logger.Logger.Infow("AlertHandler: key=%s, value=%s", k, v)
	return nil
}

var DefaultRegistry = Registry{
	"dashboard-service-my-gogin-skeleton": DashboardHandler,
	"alert-service-my-gogin-skeleton":     AlertHandler,
}
