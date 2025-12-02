package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config holds all configuration for the application.
// The values are read by viper from a config file or environment variables.

type Config struct {
	DB      Database      `mapstructure:"database"`
	Server  Server        `mapstructure:"server"`
	JWT     JWT           `mapstructure:"jwt"`
	Admin   Admin         `mapstructure:"default_admin"`
	Cache   Cache         `mapstructure:"cache"`
	Storage StorageConfig `mapstructure:"storage"`
	Email   EmailConfig   `mapstructure:"email"`
}

type Cache struct {
	Type      string          `mapstructure:"type"`
	Redis     RedisConfig     `mapstructure:"redis"`
	Ristretto RistrettoConfig `mapstructure:"ristretto"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type RistrettoConfig struct {
	MaxCost     int64 `mapstructure:"max_cost"`
	NumCounters int64 `mapstructure:"num_counters"`
	BufferItems int64 `mapstructure:"buffer_items"`
	Metrics     bool  `mapstructure:"metrics"`
}

type Database struct {
	Driver      string `mapstructure:"driver"`
	DSN         string `mapstructure:"dsn"`
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	Name        string `mapstructure:"name"`
	SSLMode     string `mapstructure:"sslmode"`
	AutoMigrate bool   `mapstructure:"auto_migrate"`
}

type Server struct {
	Env     string              `mapstructure:"env"`
	HTTP    HTTPServerConfig    `mapstructure:"http_server"`
	GRPC    GRPCServerConfig    `mapstructure:"grpc_server"`
	GraphQL GraphQLServerConfig `mapstructure:"graphql_server"`
}

type GraphQLServerConfig struct {
	Port          string   `mapstructure:"port"`
	Enable        bool     `mapstructure:"enable"`
	CorsOrigins   []string `mapstructure:"cors_origins"`
	ReadTimeout   string   `mapstructure:"read_timeout"`
	WriteTimeout  string   `mapstructure:"write_timeout"`
	IdleTimeout   string   `mapstructure:"idle_timeout"`
	StartupBanner bool     `mapstructure:"startup_banner"`
}

type HTTPServerConfig struct {
	Port          string   `mapstructure:"port"`
	Enable        bool     `mapstructure:"enable"`
	CorsOrigins   []string `mapstructure:"cors_origins"`
	ReadTimeout   string   `mapstructure:"read_timeout"`
	WriteTimeout  string   `mapstructure:"write_timeout"`
	IdleTimeout   string   `mapstructure:"idle_timeout"`
	StartupBanner bool     `mapstructure:"startup_banner"`
}

type GRPCServerConfig struct {
	Port                   string `mapstructure:"port"`
	Enable                 bool   `mapstructure:"enable"`
	MaxConnectionIdle      string `mapstructure:"max_connection_idle"`
	Timeout                string `mapstructure:"timeout"`
	MaxConnectionAge       string `mapstructure:"max_connection_age"`
	MaxConnectionAgeGrace  string `mapstructure:"max_connection_age_grace"`
	Time                   string `mapstructure:"time"`
	ForceTransportSecurity bool   `mapstructure:"force_transport_security"`
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

// LocalStorageConfig holds configuration for local storage.
type LocalStorageConfig struct {
	BasePath string `mapstructure:"base_path"`
}

// S3StorageConfig holds configuration for S3 storage.
type S3StorageConfig struct {
	Region          string `mapstructure:"region"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	Bucket          string `mapstructure:"bucket"`
}

// GCSStorageConfig holds configuration for GCS storage.
type GCSStorageConfig struct {
	ProjectID       string `mapstructure:"project_id"`
	Bucket          string `mapstructure:"bucket"`
	CredentialsFile string `mapstructure:"credentials_file"`
}

type StorageConfig struct {
	Type  string             `mapstructure:"type"`
	Local LocalStorageConfig `mapstructure:"local"`
	S3    S3StorageConfig    `mapstructure:"s3"`
	GCS   GCSStorageConfig   `mapstructure:"gcs"`
}

// SmtpConfig holds configuration for SMTP email client.
type SmtpConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
}

// SendGridConfig holds configuration for SendGrid email client.
type SendGridConfig struct {
	APIKey string `mapstructure:"api_key"`
	From   string `mapstructure:"from"`
}

type EmailConfig struct {
	Type     string         `mapstructure:"type"`
	SMTP     SmtpConfig     `mapstructure:"smtp"`
	SendGrid SendGridConfig `mapstructure:"sendgrid"`
}

// LoadConfig loads configuration from file or environment variables.
// It first reads config.yaml (base configuration), then merges environment-specific
// config (e.g., config.development.yaml), and finally allows environment variables
// to override any setting.
func LoadConfig(log *logrus.Logger) (*Config, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// Load default .env file
	if err := godotenv.Load(); err != nil {
		log.Warnf("Error loading default .env file: %v", err)
	}

	// Get environment from env var
	env := os.Getenv("SERVER_ENV")
	log.Infof("Configuration loaded for environment: %s", env)
	if env == "" {
		env = "development" // Default environment
	}

	// Merge environment-specific config if it exists
	viper.SetConfigName("config." + env)
	if err := viper.MergeInConfig(); err != nil {
		log.Warnf("Error reading environment-specific config file (config.%s.yaml): %v", env, err)
	}

	// Allow environment variables to override config file settings
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))
	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
