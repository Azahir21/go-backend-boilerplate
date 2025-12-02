package module

import (
	"github.com/azahir21/go-backend-boilerplate/ent"
	"github.com/azahir21/go-backend-boilerplate/infrastructure/cache"
	"github.com/azahir21/go-backend-boilerplate/infrastructure/external"
	"github.com/azahir21/go-backend-boilerplate/infrastructure/storage"
	"github.com/azahir21/go-backend-boilerplate/internal/shared/unitofwork"
	"github.com/sirupsen/logrus"
)

// Dependencies holds all shared infrastructure dependencies
// that modules can use for dependency injection.
type Dependencies struct {
	Log         *logrus.Logger
	DBClient    *ent.Client
	Cache       cache.Cache
	Storage     storage.Storage
	EmailClient external.EmailClient
	UoW         unitofwork.UnitOfWork
}
