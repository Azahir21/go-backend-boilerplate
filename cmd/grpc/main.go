package main

import (
	grpcapp "github.com/azahir21/go-backend-boilerplate/cmd/app/grpc"
	"github.com/azahir21/go-backend-boilerplate/pkg/logger"
)

func main() {
	log := logger.NewLogger()
	if err := grpcapp.Run(log); err != nil {
		log.Fatalf("gRPC application failed to start: %v", err)
	}
}
