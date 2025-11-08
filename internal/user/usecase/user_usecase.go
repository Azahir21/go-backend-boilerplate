package usecase

import (
	"context"
	"errors"

	"github.com/azahir21/go-backend-boilerplate/ent"
	"github.com/azahir21/go-backend-boilerplate/internal/shared/entity"
	"github.com/azahir21/go-backend-boilerplate/internal/shared/helper"
	"github.com/azahir21/go-backend-boilerplate/internal/shared/unitofwork"
	"github.com/azahir21/go-backend-boilerplate/internal/user/delivery/http/dto"
	"github.com/azahir21/go-backend-boilerplate/internal/user/repository"
)

var (
	errUsernameExists     = errors.New("username already exists")
	errEmailExists        = errors.New("email already exists")
	errInvalidCredentials = errors.New("invalid credentials")
)

type UserUsecase interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error)
	GetProfile(ctx context.Context, userID uint) (*entity.User, error)
}

type userUsecase struct {
	userRepo repository.UserRepository
	uow      unitofwork.UnitOfWork
}

func NewUserUsecase(userRepo repository.UserRepository, uow unitofwork.UnitOfWork) UserUsecase {
	return &userUsecase{userRepo: userRepo, uow: uow}
}

func (u *userUsecase) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	if err := u.ensureUsernameAvailable(ctx, req.Username); err != nil {
		return nil, err
	}
	if err := u.ensureEmailAvailable(ctx, req.Email); err != nil {
		return nil, err
	}

	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user",
	}

	err = u.uow.Do(ctx, func(txUow unitofwork.UnitOfWork) error {
		if err := txUow.UserRepository().Create(ctx, user); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return u.buildAuthResponse(user)
}

func (u *userUsecase) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := u.uow.UserRepository().FindByUsername(ctx, req.Username)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errInvalidCredentials
		}
		return nil, err
	}

	if err := helper.ComparePassword(req.Password, user.Password); err != nil {
		return nil, errInvalidCredentials
	}

	return u.buildAuthResponse(user)
}

func (u *userUsecase) GetProfile(ctx context.Context, userID uint) (*entity.User, error) {
	return u.uow.UserRepository().FindByID(ctx, userID)
}

func (u *userUsecase) ensureUsernameAvailable(ctx context.Context, username string) error {
	_, err := u.uow.UserRepository().FindByUsername(ctx, username)
	if err == nil {
		return errUsernameExists
	}
	if ent.IsNotFound(err) {
		return nil
	}
	return err
}

func (u *userUsecase) ensureEmailAvailable(ctx context.Context, email string) error {
	_, err := u.uow.UserRepository().FindByEmail(ctx, email)
	if err == nil {
		return errEmailExists
	}
	if ent.IsNotFound(err) {
		return nil
	}
	return err
}

func (u *userUsecase) buildAuthResponse(user *entity.User) (*dto.AuthResponse, error) {
	token, err := helper.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token: token,
		User:  dto.ToUserResponse(user),
	}, nil
}
