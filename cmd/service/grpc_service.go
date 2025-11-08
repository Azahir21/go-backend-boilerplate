package service

import (
	"fmt"
	"time"

	grpcDelivery "github.com/azahir21/go-backend-boilerplate/internal/user/delivery/grpc"
	userUsecase "github.com/azahir21/go-backend-boilerplate/internal/user/usecase"
	"github.com/azahir21/go-backend-boilerplate/pkg/config"
	"github.com/azahir21/go-backend-boilerplate/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func NewGrpcServer(log *logrus.Logger, cfg config.GRPCServerConfig, userUsecase userUsecase.UserUsecase) (*grpc.Server, error) {
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

	userGrpcHandler := grpcDelivery.NewUserHandler(log, userUsecase)
	proto.RegisterUserServiceServer(grpcServer, userGrpcHandler)

	if cfg.Enable {
		log.Info("gRPC server configured successfully")
	}

	return grpcServer, nil
}
