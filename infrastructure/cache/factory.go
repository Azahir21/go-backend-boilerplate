package cache

import (
	"fmt"

	"github.com/azahir21/go-backend-boilerplate/pkg/config"
	"github.com/sirupsen/logrus"
)

func NewCache(log *logrus.Logger, cfg config.Cache) (Cache, error) {
	switch cfg.Type {
	case "redis":
		return NewRedisCache(log, cfg.Redis)
	case "ristretto":
		return NewRistrettoCache(log, cfg.Ristretto)
	default:
		return nil, fmt.Errorf("unsupported cache type: %s", cfg.Type)
	}
}
