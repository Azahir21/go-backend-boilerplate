//go:build graphql
// +build graphql

package graphql

import (
	"context"
	"errors"

	"github.com/azahir21/go-backend-boilerplate/internal/user/delivery/http/dto"
	"github.com/azahir21/go-backend-boilerplate/internal/user/usecase"
	"github.com/graphql-go/graphql"
	"github.com/sirupsen/logrus"
)

type UserResolver struct {
	log         *logrus.Logger
	userUsecase usecase.UserUsecase
}

func NewUserResolver(log *logrus.Logger, userUsecase usecase.UserUsecase) *UserResolver {
	return &UserResolver{
		log:         log,
		userUsecase: userUsecase,
	}
}

func (r *UserResolver) GetUserResolver(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(int)
	if !ok {
		return nil, errors.New("invalid user ID")
	}

	// Use context from GraphQL params if available
	ctx := p.Context
	if ctx == nil {
		ctx = context.Background()
	}

	user, err := r.userUsecase.GetProfile(ctx, uint(id))
	if err != nil {
		r.log.Errorf("failed to get user profile for ID %d: %v", id, err)
		return nil, err
	}

	return user, nil
}

func (r *UserResolver) RegisterUserResolver(p graphql.ResolveParams) (interface{}, error) {
	username, usernameOk := p.Args["username"].(string)
	email, emailOk := p.Args["email"].(string)
	password, passwordOk := p.Args["password"].(string)

	if !usernameOk || !emailOk || !passwordOk {
		return nil, errors.New("invalid registration input")
	}

	req := &dto.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	}

	// Use context from GraphQL params if available
	ctx := p.Context
	if ctx == nil {
		ctx = context.Background()
	}

	authResponse, err := r.userUsecase.Register(ctx, req)
	if err != nil {
		r.log.Errorf("failed to register user %s: %v", username, err)
		return nil, err
	}

	return authResponse.User, nil
}

func (r *UserResolver) LoginUserResolver(p graphql.ResolveParams) (interface{}, error) {
	username, usernameOk := p.Args["username"].(string)
	password, passwordOk := p.Args["password"].(string)

	if !usernameOk || !passwordOk {
		return nil, errors.New("invalid login input")
	}

	req := &dto.LoginRequest{
		Username: username,
		Password: password,
	}

	// Use context from GraphQL params if available
	ctx := p.Context
	if ctx == nil {
		ctx = context.Background()
	}

	authResponse, err := r.userUsecase.Login(ctx, req)
	if err != nil {
		r.log.Errorf("failed to login user %s: %v", username, err)
		return nil, err
	}

	return authResponse.Token, nil
}
