package base

import (
	"context"
	"time"
)

// QueryContext : The context used for db queries
func QueryContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}
