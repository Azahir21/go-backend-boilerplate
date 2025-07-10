package service

import (
	"errors"

	"github.com/azahir21/go-backend-boilerplate/internal/helper"
	"github.com/azahir21/go-backend-boilerplate/internal/model"
	"github.com/azahir21/go-backend-boilerplate/internal/repository"
	"gorm.io/gorm"
)

type AuthService interface {
    Register(req *model.RegisterRequest) (*model.AuthResponse, error)
    Login(req *model.LoginRequest) (*model.AuthResponse, error)
    GetProfile(userID uint) (*model.User, error)
}

type authService struct {
    userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
    return &authService{userRepo: userRepo}
}

func (s *authService) Register(req *model.RegisterRequest) (*model.AuthResponse, error) {
    // Check if user already exists
    if _, err := s.userRepo.FindByUsername(req.Username); err == nil {
        return nil, errors.New("username already exists")
    }

    if _, err := s.userRepo.FindByEmail(req.Email); err == nil {
        return nil, errors.New("email already exists")
    }

    // Hash password
    hashedPassword, err := helper.HashPassword(req.Password)
    if err != nil {
        return nil, err
    }

    // Create user
    user := &model.User{
        Username: req.Username,
        Email:    req.Email,
        Password: hashedPassword,
        Role:     "user",
    }

    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }

    // Generate token
    token, err := helper.GenerateToken(user.ID, user.Username, user.Role)
    if err != nil {
        return nil, err
    }

    return &model.AuthResponse{
        Token: token,
        User:  user.ToUserResponse(),
    }, nil
}

func (s *authService) Login(req *model.LoginRequest) (*model.AuthResponse, error) {
    // Find user by username
    user, err := s.userRepo.FindByUsername(req.Username)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("invalid credentials")
        }
        return nil, err
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

    return &model.AuthResponse{
        Token: token,
        User:  user.ToUserResponse(),
    }, nil
}

func (s *authService) GetProfile(userID uint) (*model.User, error) {
    return s.userRepo.FindByID(userID)
}