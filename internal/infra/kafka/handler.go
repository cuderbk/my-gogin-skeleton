package kafka

import "context"

/* Signature for all handlers */
type Handler func(ctx context.Context, key, value []byte) error

/* Topic → Handler map */
type Registry map[string]Handler
