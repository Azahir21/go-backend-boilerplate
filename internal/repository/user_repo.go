package repository

import (
	"context"

	"github.com/azahir21/go-backend-boilerplate/internal/domain"
)

type UserRepository interface {
    Create(ctx context.Context, user *domain.User) error
    FindByUsername(ctx context.Context, username string) (*domain.User, error)
    FindByEmail(ctx context.Context, email string) (*domain.User, error)
    FindByID(ctx context.Context, id uint) (*domain.User, error)
    Update(ctx context.Context, user *domain.User) error
    Delete(ctx context.Context, id uint) error
}
