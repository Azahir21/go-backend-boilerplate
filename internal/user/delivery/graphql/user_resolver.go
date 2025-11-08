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

	user, err := r.userUsecase.GetProfile(context.Background(), uint(id))
	if err != nil {
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

	authResponse, err := r.userUsecase.Register(context.Background(), req)
	if err != nil {
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

	authResponse, err := r.userUsecase.Login(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return authResponse.Token, nil
}
