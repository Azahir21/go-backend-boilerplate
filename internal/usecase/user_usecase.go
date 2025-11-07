package usecase

import (
	"context"
	"errors"

	"github.com/azahir21/go-backend-boilerplate/internal/domain"
	"github.com/azahir21/go-backend-boilerplate/internal/helper"
	"github.com/azahir21/go-backend-boilerplate/internal/repository"
)

type UserUsecase interface {
    Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error)
    Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error)
    GetProfile(ctx context.Context, userID uint) (*domain.User, error)
}

type userUsecase struct {
    userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
    return &userUsecase{userRepo: userRepo}
}

func (u *userUsecase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error) {
    // Check if user already exists
    if _, err := u.userRepo.FindByUsername(ctx, req.Username); err == nil {
        return nil, errors.New("username already exists")
    }

    if _, err := u.userRepo.FindByEmail(ctx, req.Email); err == nil {
        return nil, errors.New("email already exists")
    }

    // Hash password
    hashedPassword, err := helper.HashPassword(req.Password)
    if err != nil {
        return nil, err
    }

    // Create user
    user := &domain.User{
        Username: req.Username,
        Email:    req.Email,
        Password: hashedPassword,
        Role:     "user",
    }

    if err := u.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }

    // Generate token
    token, err := helper.GenerateToken(user.ID, user.Username, user.Role)
    if err != nil {
        return nil, err
    }

    return &domain.AuthResponse{
        Token: token,
        User:  user.ToUserResponse(),
    }, nil
}

func (u *userUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
    // Find user by username
    user, err := u.userRepo.FindByUsername(ctx, req.Username)
    if err != nil {
        return nil, errors.New("invalid credentials")
    }

    // Compare password
    if err := helper.ComparePassword(req.Password, user.Password); err != nil {
        return nil, errors.New("invalid credentials")
    }

    // Generate token
    token, err := helper.GenerateToken(user.ID, user.Username, user.Role)
    if err != nil {
        return nil, err
    }

    return &domain.AuthResponse{
        Token: token,
        User:  user.ToUserResponse(),
    }, nil
}

func (u *userUsecase) GetProfile(ctx context.Context, userID uint) (*domain.User, error) {
    return u.userRepo.FindByID(ctx, userID)
}
