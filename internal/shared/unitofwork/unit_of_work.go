package unitofwork

import (
	"context"
	"fmt"

	"github.com/azahir21/go-backend-boilerplate/ent"
	"github.com/azahir21/go-backend-boilerplate/internal/user/repository"
	userRepoImpl "github.com/azahir21/go-backend-boilerplate/internal/user/repository/implementation"
)

// UnitOfWork defines the interface for a unit of work.
type UnitOfWork interface {
	Do(ctx context.Context, fn func(txUow UnitOfWork) error) error
	UserRepository() repository.UserRepository
	// Add other repositories here
}

type unitOfWork struct {
	client *ent.Client
	tx     *ent.Tx
}

// NewUnitOfWork creates a new UnitOfWork instance.
func NewUnitOfWork(client *ent.Client) UnitOfWork {
	return &unitOfWork{client: client}
}

// Do executes a function within a database transaction.
func (u *unitOfWork) Do(ctx context.Context, fn func(txUow UnitOfWork) error) error {
	tx, err := u.client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()

	txUow := &unitOfWork{client: tx.Client(), tx: tx}
	if err := fn(txUow); err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// UserRepository returns a transactional UserRepository.
func (u *unitOfWork) UserRepository() repository.UserRepository {
	return userRepoImpl.NewUserRepository(u.client)
}

// Add other transactional repository getters here
