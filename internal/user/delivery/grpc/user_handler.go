package grpc

import (
	"context"

	"github.com/azahir21/go-backend-boilerplate/internal/user/usecase"
	"github.com/azahir21/go-backend-boilerplate/proto"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	proto.UnimplementedUserServiceServer
	log         *logrus.Logger
	userUsecase usecase.UserUsecase
}

func NewUserHandler(log *logrus.Logger, userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		log:         log,
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.AuthResponse, error) {
	// Placeholder implementation
	return &proto.AuthResponse{Token: "dummy-token"}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *proto.LoginRequest) (*proto.AuthResponse, error) {
	// Placeholder implementation
	return &proto.AuthResponse{Token: "dummy-token"}, nil
}

func (h *UserHandler) GetProfile(ctx context.Context, req *proto.GetProfileRequest) (*proto.User, error) {
	// Placeholder implementation
	return &proto.User{Id: req.UserId, Username: "dummy-user"}, nil
}
