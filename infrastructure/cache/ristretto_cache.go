package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/azahir21/go-backend-boilerplate/pkg/config"
	"github.com/dgraph-io/ristretto"
	"github.com/sirupsen/logrus"
)

type RistrettoCache struct {
	cached *ristretto.Cache
	log    *logrus.Logger
}

func NewRistrettoCache(log *logrus.Logger, cfg config.RistrettoConfig) (*RistrettoCache, error) {
	cached, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: cfg.NumCounters,
		MaxCost:     cfg.MaxCost,
		BufferItems: cfg.BufferItems,
		Metrics:     cfg.Metrics,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize ristretto cache: %w", err)
	}

	log.Info("Ristretto in-memory cache initialized")
	return &RistrettoCache{cached: cached, log: log}, nil
}

func (r *RistrettoCache) Set(ctx context.Context, key string, value interface{}, ttlSeconds int) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value for key %s: %w", key, err)
	}

	if !r.cached.SetWithTTL(key, data, 1, time.Duration(ttlSeconds)*time.Second) {
		return fmt.Errorf("failed to set cache key %s: cache rejected the entry", key)
	}
	return nil
}

func (r *RistrettoCache) Get(ctx context.Context, key string, dest interface{}) error {
	value, found := r.cached.Get(key)
	if !found {
		return fmt.Errorf("cache key not found: %s", key)
	}

	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cache value for key %s has unexpected type: %T", key, value)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal cache value for key %s: %w", key, err)
	}
	return nil
}

func (r *RistrettoCache) Del(ctx context.Context, key string) error {
	r.cached.Del(key)
	return nil
}
