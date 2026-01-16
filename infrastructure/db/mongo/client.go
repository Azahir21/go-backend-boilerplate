package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/azahir21/go-backend-boilerplate/pkg/config"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Client wraps the MongoDB client and database handle.
type Client struct {
	client   *mongo.Client
	database *mongo.Database
	log      *logrus.Logger
}

// NewClient creates a new MongoDB client and establishes a connection.
func NewClient(log *logrus.Logger, cfg config.MongoConfig) (*Client, error) {
	// Build connection options
	clientOpts := options.Client().ApplyURI(cfg.URI)

	// Add authentication if credentials are provided
	if cfg.Username != "" && cfg.Password != "" {
		authSource := cfg.AuthSource
		if authSource == "" {
			authSource = "admin" // Default auth source
		}
		credential := options.Credential{
			Username:   cfg.Username,
			Password:   cfg.Password,
			AuthSource: authSource,
		}
		clientOpts.SetAuth(credential)
	}

	// Set connection timeout
	connectTimeout := time.Duration(cfg.ConnectTimeoutMS) * time.Millisecond
	if connectTimeout == 0 {
		connectTimeout = 10 * time.Second // Default timeout
	}
	clientOpts.SetConnectTimeout(connectTimeout)

	// Set pool size if configured
	if cfg.MaxPoolSize > 0 {
		clientOpts.SetMaxPoolSize(cfg.MaxPoolSize)
	}
	if cfg.MinPoolSize > 0 {
		clientOpts.SetMinPoolSize(cfg.MinPoolSize)
	}

	// Create context with timeout for connection
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Verify connection with a ping
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer pingCancel()

	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		client.Disconnect(context.Background())
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	database := client.Database(cfg.Database)
	log.Infof("MongoDB connection established (database: %s, uri: %s)", cfg.Database, cfg.URI)

	return &Client{
		client:   client,
		database: database,
		log:      log,
	}, nil
}

// Database returns the MongoDB database handle.
func (c *Client) Database() *mongo.Database {
	return c.database
}

// Client returns the underlying MongoDB client.
func (c *Client) Client() *mongo.Client {
	return c.client
}

// Close closes the MongoDB connection.
func (c *Client) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := c.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to close MongoDB connection: %w", err)
	}

	c.log.Info("MongoDB connection closed")
	return nil
}

// Collection returns a handle to the specified collection.
func (c *Client) Collection(name string) *mongo.Collection {
	return c.database.Collection(name)
}
