package cache

import (
	"context"
	"encoding/json"
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

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	log.Info("Redis cache initialized")
	return &RedisCache{client: client, log: log}, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, ttlSeconds int) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func (r *RedisCache) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
