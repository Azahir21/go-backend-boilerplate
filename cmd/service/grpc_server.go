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

func NewGrpcServer(log *logrus.Logger, cfg config.GRPCServerConfig, userUsecase userUsecase.UserUsecase) *grpc.Server {

	maxConnectionIdle, err := time.ParseDuration(cfg.MaxConnectionIdle)
	if err != nil {
		log.Fatalf("invalid max connection idle duration: %v", err)
	}

	timeDuration, err := time.ParseDuration(cfg.Time)
	if err != nil {
		log.Fatalf("invalid time duration: %v", err)
	}

	timeout, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		log.Fatalf("invalid timeout duration: %v", err)
	}

	maxConnectionAge, err := time.ParseDuration(cfg.MaxConnectionAge)
	if err != nil {
		log.Fatalf("invalid max connection age duration: %v", err)
	}

	maxConnectionAgeGrace, err := time.ParseDuration(cfg.MaxConnectionAgeGrace)
	if err != nil {
		log.Fatalf("invalid max connection age grace duration: %v", err)
	}


	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: maxConnectionIdle,
			Time:              timeDuration,
			Timeout:           timeout,
			MaxConnectionAge:  maxConnectionAge,
			MaxConnectionAgeGrace: maxConnectionAgeGrace,
		}),
	)

	userGrpcHandler := grpcDelivery.NewUserHandler(log, userUsecase)
	proto.RegisterUserServiceServer(grpcServer, userGrpcHandler)

	if cfg.Enable {
		fmt.Println("🚀 Starting gRPC server...")
	}

	return grpcServer
}
