package unitofwork

import (
	"context"
	"fmt"

	"github.com/azahir21/go-backend-boilerplate/ent"
	"github.com/azahir21/go-backend-boilerplate/internal/user/repository"
	userRepoImpl "github.com/azahir21/go-backend-boilerplate/internal/user/repository/implementation"
)

// UnitOfWork defines the interface for a unit of work pattern.
// It manages database transactions and provides access to repositories
// within the transaction boundary.
type UnitOfWork interface {
	// Do executes a function within a database transaction.
	// If the function returns an error, the transaction is rolled back.
	// If the function succeeds, the transaction is committed.
	Do(ctx context.Context, fn func(txUow UnitOfWork) error) error

	// UserRepository returns the user repository for this unit of work.
	// If called within a transaction, it returns a transactional repository.
	UserRepository() repository.UserRepository
	// Add other repository getters here as needed
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
// It handles automatic rollback on errors and panic recovery.
func (u *unitOfWork) Do(ctx context.Context, fn func(txUow UnitOfWork) error) error {
	tx, err := u.client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Ensure transaction is rolled back on panic
	defer func() {
		if v := recover(); v != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				// Log rollback error but re-panic with original
				fmt.Printf("failed to rollback transaction after panic: %v\n", rollbackErr)
			}
			panic(v)
		}
	}()

	// Create transactional unit of work
	txUow := &unitOfWork{client: tx.Client(), tx: tx}

	// Execute the function
	if err := fn(txUow); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("transaction failed (rollback also failed: %v): %w", rollbackErr, err)
		}
		return fmt.Errorf("transaction failed: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// UserRepository returns a repository instance.
// If this unit of work is transactional, the repository will use the transaction.
func (u *unitOfWork) UserRepository() repository.UserRepository {
	return userRepoImpl.NewUserRepository(u.client)
}
