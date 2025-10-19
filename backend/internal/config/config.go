package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	Server        ServerConfig
	JWT           JWTConfig
	Unipile       UnipileConfig
	JWTSecret     string
	UnipileAPIKey string
	DatabasePath  string
	FrontendURL   string
}

type ServerConfig struct {
	Port int
}

type JWTConfig struct {
	TokenDuration time.Duration `mapstructure:"token_duration"`
}

type UnipileConfig struct {
	APIURL        string
	Timeout       time.Duration
	RetryAttempts int           `mapstructure:"retry_attempts"`
	RetryDelay    time.Duration `mapstructure:"retry_delay"`
}

var App *Config

// LoadConfig loads configuration from YAML and environment variables
func LoadConfig() error {
	// Load .env file for secrets (don't fail if not found)
	if err := godotenv.Load("configs/.env"); err != nil {
		// Try root directory
		godotenv.Load(".env")
	}

	// Initialize Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../configs")
	viper.AddConfigPath(".")

	// Read base config
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	log.Printf("Loaded config from: %s", viper.ConfigFileUsed())

	// Check for environment-specific config
	env := getEnv("APP_ENV", viper.GetString("app.environment"))

	if env == "production" || env == "staging" {
		viper.SetConfigName("config." + env)
		if err := viper.MergeInConfig(); err != nil {
			log.Printf("No %s config found, using defaults", env)
		} else {
			log.Printf("Loaded %s overrides", env)
		}
	}

	// Unmarshal into struct
	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Load secrets from environment variables
	cfg.JWTSecret = getEnv("JWT_SECRET", "")
	if cfg.JWTSecret == "" {
		log.Println("WARNING: JWT_SECRET not set! Using insecure default.")
		cfg.JWTSecret = "default-secret-change-in-production"
	}

	cfg.UnipileAPIKey = getEnv("UNIPILE_API_KEY", "")
	if cfg.UnipileAPIKey == "" {
		log.Println("WARNING: UNIPILE_API_KEY not set!")
	}

	cfg.DatabasePath = getEnv("DATABASE_PATH", "./linkedin_connector.db")
	cfg.FrontendURL = getEnv("FRONTEND_URL", "http://localhost:5173")

	// Set Unipile API URL (can be overridden by env)
	cfg.Unipile.APIURL = getEnv("UNIPILE_API_URL", "https://api.unipile.com/v1")

	// Override server port from env if provided (for deployment platforms)
	if portStr := os.Getenv("PORT"); portStr != "" {
		var port int
		if _, err := fmt.Sscanf(portStr, "%d", &port); err == nil {
			cfg.Server.Port = port
		}
	}

	// Assign to global App variable
	App = cfg

	return nil
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetConfig() *Config {
	return App
}
