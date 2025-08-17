package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config represents the top-level application configuration
// loaded from a YAML file and environment variables.
type Config struct {
	Env      string   `yaml:"env"`
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
	Kafka    Kafka    `yaml:"kafka"`
	Cache    Cache    `yaml:"cache"`
}

// Server holds configuration for the HTTP server.
type Server struct {
	HTTPPort string `yaml:"httpPort"`
}

// Database holds database connection configuration.
// Some fields (host, port, user, password, name) are
// expected to be overridden from environment variables.
type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string `yaml:"sslmode"`
}

// Kafka holds configuration for Kafka consumers/producers.
type Kafka struct {
	GroupID string   `yaml:"groupID"`
	Topic   string   `yaml:"topic"`
	Brokers []string `yaml:"brokers"`
}

// Cache holds configuration for in-memory cache.
type Cache struct {
	DefaultExpiration time.Duration `yaml:"defaultExpiration"`
	CleanupInterval   time.Duration `yaml:"cleanupInterval"`
	PreloadLimit      int           `yaml:"preloadLimit"`
}

// DatabaseURL builds a PostgreSQL connection string
// based on the Database configuration.
func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// MustLoad loads application configuration from the ./config/config.yaml file
// and environment variables. If the configuration file cannot be read or parsed,
// the function logs a fatal error and terminates the program.
//
// Values from environment variables (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
// override the corresponding fields in the Database section of the config file.
func MustLoad() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("fatal error config file: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	cfg.Database.Host = os.Getenv("DB_HOST")
	cfg.Database.Port = os.Getenv("DB_PORT")
	cfg.Database.User = os.Getenv("DB_USER")
	cfg.Database.Password = os.Getenv("DB_PASSWORD")
	cfg.Database.Name = os.Getenv("DB_NAME")

	return &cfg
}
