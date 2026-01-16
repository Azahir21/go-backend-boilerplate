package service

import (
	"fmt"
	"time"

	"github.com/azahir21/go-backend-boilerplate/internal/shared/module"
	"github.com/azahir21/go-backend-boilerplate/pkg/config"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func NewGrpcServer(log *logrus.Logger, cfg config.GRPCServerConfig, modules []module.GRPCModule) (*grpc.Server, error) {
	maxConnectionIdle, err := time.ParseDuration(cfg.MaxConnectionIdle)
	if err != nil {
		return nil, fmt.Errorf("invalid max connection idle duration: %w", err)
	}

	timeDuration, err := time.ParseDuration(cfg.Time)
	if err != nil {
		return nil, fmt.Errorf("invalid time duration: %w", err)
	}

	timeout, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		return nil, fmt.Errorf("invalid timeout duration: %w", err)
	}

	maxConnectionAge, err := time.ParseDuration(cfg.MaxConnectionAge)
	if err != nil {
		return nil, fmt.Errorf("invalid max connection age duration: %w", err)
	}

	maxConnectionAgeGrace, err := time.ParseDuration(cfg.MaxConnectionAgeGrace)
	if err != nil {
		return nil, fmt.Errorf("invalid max connection age grace duration: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     maxConnectionIdle,
			Time:                  timeDuration,
			Timeout:               timeout,
			MaxConnectionAge:      maxConnectionAge,
			MaxConnectionAgeGrace: maxConnectionAgeGrace,
		}),
	)

	// Register gRPC services from all modules
	for _, m := range modules {
		log.Infof("Registering gRPC services for module: %s", m.Name())
		m.RegisterGRPC(grpcServer)
	}

	if cfg.Enable {
		log.Infof("🚀 Starting gRPC server on :%s", cfg.Port)
	}

	return grpcServer, nil
}
