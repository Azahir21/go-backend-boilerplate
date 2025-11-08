package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/azahir21/go-backend-boilerplate/pkg/config"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisCache struct {
	client *redis.Client
	log    *logrus.Logger
}

func NewRedisCache(log *logrus.Logger, cfg config.RedisConfig) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis at %s: %w", cfg.Addr, err)
	}

	log.Infof("Redis cache initialized at %s", cfg.Addr)
	return &RedisCache{client: client, log: log}, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttlSeconds int) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value for key %s: %w", key, err)
	}

	if err := r.client.Set(ctx, key, data, time.Duration(ttlSeconds)*time.Second).Err(); err != nil {
		return fmt.Errorf("failed to set cache key %s: %w", key, err)
	}
	return nil
}

func (r *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("cache key not found: %s", key)
		}
		return fmt.Errorf("failed to get cache key %s: %w", key, err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal cache value for key %s: %w", key, err)
	}
	return nil
}

func (r *RedisCache) Del(ctx context.Context, key string) error {
	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete cache key %s: %w", key, err)
	}
	return nil
}
