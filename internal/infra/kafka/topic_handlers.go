package kafka

import (
	"context"
	"logging/internal/common/logger"
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
	"dashboard-service-logging": DashboardHandler,
	"alert-service-logging":     AlertHandler,
}
