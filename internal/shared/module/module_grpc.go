//go:build grpc
// +build grpc

package module

import (
	"google.golang.org/grpc"
)

// GRPCModule is implemented by modules that provide gRPC handlers.
type GRPCModule interface {
	Module
	// RegisterGRPC registers the module's gRPC services on the given server.
	RegisterGRPC(server *grpc.Server)
}
