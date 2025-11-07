package implementation

import (
	"context"

	"github.com/azahir21/go-backend-boilerplate/ent"
	"github.com/azahir21/go-backend-boilerplate/ent/user"
	"github.com/azahir21/go-backend-boilerplate/internal/shared/entity"
	"github.com/azahir21/go-backend-boilerplate/internal/user/repository"
)

type userRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) repository.UserRepository {
	return &userRepository{client: client}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	_, err := r.client.User.
		Create().
		SetUsername(user.Username).
		SetEmail(user.Email).
		SetPassword(user.Password).
		SetRole(user.Role).
		Save(ctx)
	return err
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	entUser, err := r.client.User.Query().Where(user.UsernameEQ(username)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return toDomainUser(entUser), nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	entUser, err := r.client.User.Query().Where(user.EmailEQ(email)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return toDomainUser(entUser), nil
}

func (r *userRepository) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	entUser, err := r.client.User.Get(ctx, int(id))
	if err != nil {
		return nil, err
	}
	return toDomainUser(entUser), nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	_, err := r.client.User.
		UpdateOneID(int(user.ID)).
		SetUsername(user.Username).
		SetEmail(user.Email).
		SetPassword(user.Password).
		SetRole(user.Role).
		Save(ctx)
	return err
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.client.User.DeleteOneID(int(id)).Exec(ctx)
}

func toDomainUser(entUser *ent.User) *entity.User {
	return &entity.User{
		ID:        uint(entUser.ID),
		Username:  entUser.Username,
		Email:     entUser.Email,
		Password:  entUser.Password,
		Role:      entUser.Role,
		CreatedAt: entUser.CreatedAt,
		UpdatedAt: entUser.UpdatedAt,
		DeletedAt: entUser.DeletedAt,
	}
}
