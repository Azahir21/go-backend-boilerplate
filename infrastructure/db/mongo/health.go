package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// HealthCheck performs a health check on the MongoDB connection.
// It returns nil if the connection is healthy, otherwise returns an error.
func (c *Client) HealthCheck(ctx context.Context) error {
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return c.client.Ping(pingCtx, readpref.Primary())
}

// IsHealthy returns true if the MongoDB connection is healthy.
func (c *Client) IsHealthy() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.HealthCheck(ctx) == nil
}
