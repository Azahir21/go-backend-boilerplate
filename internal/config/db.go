package config

import (
	"fmt"

	"github.com/azahir21/go-backend-boilerplate/internal/helper"
	"github.com/azahir21/go-backend-boilerplate/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(log *logrus.Logger, cfg *Config) *gorm.DB {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", 
        cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode,
    )
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect database: %v", err)
    }
    
    return db
}

func Migrate(db *gorm.DB, log *logrus.Logger, cfg *Config) {
    // Auto-migrate the User model
    err := db.AutoMigrate(&model.User{})
    if err != nil {
        log.Fatalf("failed to migrate database: %v", err)
    }
    
    // Create default admin user if not exists
    var adminUser model.User
    result := db.Where("username = ?", cfg.DefaultAdminUsername).First(&adminUser)
    if result.Error != nil {
        // Create default admin user
        hashedPassword, _ := helper.HashPassword(cfg.DefaultAdminPassword)
        admin := model.User{
            Username: cfg.DefaultAdminUsername,
            Email:    cfg.DefaultAdminEmail,
            Password: hashedPassword,
            Role:     "admin",
        }
        db.Create(&admin)
        log.Infof("Default admin user created (username: %s, password: %s)", 
            cfg.DefaultAdminUsername, cfg.DefaultAdminPassword)
    }
    
    log.Info("Migration completed")
}