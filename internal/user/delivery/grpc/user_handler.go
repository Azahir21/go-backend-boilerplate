package grpc

import (
	"context"
	"errors"

	"github.com/azahir21/go-backend-boilerplate/internal/user/delivery/http/dto"
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
	registerReq := &dto.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	authResponse, err := h.userUsecase.Register(ctx, registerReq)
	if err != nil {
		h.log.Errorf("gRPC Register failed: %v", err)
		return nil, errors.New("registration failed")
	}

	return &proto.AuthResponse{
		Token: authResponse.Token,
		User: &proto.User{
			Id:       uint32(authResponse.User.ID),
			Username: authResponse.User.Username,
			Email:    authResponse.User.Email,
			Role:     authResponse.User.Role,
		},
	}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *proto.LoginRequest) (*proto.AuthResponse, error) {
	loginReq := &dto.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}

	authResponse, err := h.userUsecase.Login(ctx, loginReq)
	if err != nil {
		h.log.Errorf("gRPC Login failed: %v", err)
		return nil, errors.New("login failed")
	}

	return &proto.AuthResponse{
		Token: authResponse.Token,
		User: &proto.User{
			Id:       uint32(authResponse.User.ID),
			Username: authResponse.User.Username,
			Email:    authResponse.User.Email,
			Role:     authResponse.User.Role,
		},
	}, nil
}

func (h *UserHandler) GetProfile(ctx context.Context, req *proto.GetProfileRequest) (*proto.User, error) {
	user, err := h.userUsecase.GetProfile(ctx, uint(req.UserId))
	if err != nil {
		h.log.Errorf("gRPC GetProfile failed: %v", err)
		return nil, errors.New("user not found")
	}

	return &proto.User{
		Id:       uint32(user.ID),
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}
