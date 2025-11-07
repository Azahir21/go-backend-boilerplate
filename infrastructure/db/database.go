package db

import (
	"context"
	"fmt"

	"github.com/azahir21/go-backend-boilerplate/ent"
	"github.com/azahir21/go-backend-boilerplate/pkg/config"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	_ "github.com/lib/pq"              // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3"    // SQLite driver
	"github.com/sirupsen/logrus"
)

func NewEntClient(log *logrus.Logger, cfg *config.Config) (*ent.Client, error) {
	var dsn string

	if cfg.DB.DSN != "" {
		dsn = cfg.DB.DSN
	} else {
		switch cfg.DB.Driver {
		case "postgres":
			dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
				cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLMode)
		case "mysql", "mariadb":
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True",
				cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)
		case "sqlite3":
			dsn = fmt.Sprintf("file:%s?_fk=1", cfg.DB.Name)
		default:
			return nil, fmt.Errorf("unsupported database driver: %s", cfg.DB.Driver)
		}
	}

	client, err := ent.Open(cfg.DB.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to %s: %v", cfg.DB.Driver, err)
	}

	if cfg.DB.AutoMigrate {
		// Run the auto migration tool.
		if err := client.Schema.Create(context.Background()); err != nil {
			return nil, fmt.Errorf("failed creating schema resources: %v", err)
		}
		log.Info("Ent auto-migration completed.")
	}

	log.Infof("Database connection established with driver: %s", cfg.DB.Driver)
	return client, nil
}
