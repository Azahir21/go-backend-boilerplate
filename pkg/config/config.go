package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config holds all configuration for the application.
// The values are read by viper from a config file or environment variables.

type Config struct {
	DB     Database `mapstructure:"database"`
	Server Server   `mapstructure:"server"`
	JWT    JWT      `mapstructure:"jwt"`
	Admin  Admin    `mapstructure:"default_admin"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
}

type Server struct {
	HTTPPort string `mapstructure:"http_port"`
	GRPCPort string `mapstructure:"grpc_port"`
	Env      string `mapstructure:"env"`
}

type JWT struct {
	Secret      string `mapstructure:"secret"`
	ExpiryHours int    `mapstructure:"expiry_hours"`
}

type Admin struct {
	Username string `mapstructure:"username"`
	Email    string `mapstructure:"email"`
	Password string `mapstructure:"password"`
}

// LoadConfig loads configuration from file or environment variables.
func LoadConfig(log *logrus.Logger) (*Config, error) {
	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Warn("Error reading config file, using environment variables")
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
