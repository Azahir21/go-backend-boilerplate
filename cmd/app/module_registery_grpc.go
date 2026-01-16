//go:build grpc
// +build grpc

package app

import (
	"github.com/azahir21/go-backend-boilerplate/internal/shared/module"
	userConfig "github.com/azahir21/go-backend-boilerplate/internal/user/config"
)

// registerGRPCModules registers gRPC modules when the grpc build tag is enabled.
func registerGRPCModules(deps *module.Dependencies) []module.GRPCModule {
	var grpcModules []module.GRPCModule
	
	// User module
	grpcModules = append(grpcModules, userConfig.NewGRPCConfig(deps))

	// Add more modules here as needed:
	// grpcModules = append(grpcModules, productConfig.NewGRPCConfig(deps))

	return grpcModules
}
