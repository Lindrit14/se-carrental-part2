package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	RabbitMQURL string
	HTTPPort    string

	NotifierType string // "mock" or "acs"
	FromName     string

	// ACS-specific (used when NotifierType == "acs")
	ACSEndpoint              string
	ACSSenderAddress         string
	ACSAuthMode              string // "managed_identity" or "connection_string"
	ACSConnectionStringFile  string
	ACSConnectionStringInline string

	// Redis-backed user→email read model
	RedisAddr      string
	RedisPassword  string
	RedisDB        int
	RedisKeyPrefix string

	// Used to build the password-reset link in emails.
	FrontendBaseURL string
}

func Load() (*Config, error) {
	c := &Config{
		RabbitMQURL:               getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		HTTPPort:                  getEnv("HTTP_PORT", "8084"),
		NotifierType:              getEnv("NOTIFIER_TYPE", "mock"),
		FromName:                  getEnv("FROM_NAME", "Car Rental"),
		ACSEndpoint:               getEnv("ACS_ENDPOINT", ""),
		ACSSenderAddress:          getEnv("ACS_SENDER_ADDRESS", ""),
		ACSAuthMode:               getEnv("ACS_AUTH_MODE", "managed_identity"),
		ACSConnectionStringFile:   getEnv("ACS_CONNECTION_STRING_FILE", ""),
		ACSConnectionStringInline: getEnv("ACS_CONNECTION_STRING", ""),
		RedisAddr:                 getEnv("REDIS_ADDR", "redis:6379"),
		RedisPassword:             getEnv("REDIS_PASSWORD", ""),
		RedisDB:                   getEnvInt("REDIS_DB", 0),
		RedisKeyPrefix:            getEnv("REDIS_KEY_PREFIX", "notif"),
		FrontendBaseURL:           getEnv("FRONTEND_BASE_URL", ""),
	}
	if err := c.validate(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) validate() error {
	if c.RabbitMQURL == "" {
		return fmt.Errorf("RABBITMQ_URL is required")
	}
	if _, err := strconv.Atoi(c.HTTPPort); err != nil {
		return fmt.Errorf("HTTP_PORT must be numeric: %w", err)
	}
	if c.RedisAddr == "" {
		return fmt.Errorf("REDIS_ADDR is required")
	}

	switch strings.ToLower(c.NotifierType) {
	case "mock":
		// Mock has no extra requirements.
	case "acs":
		if c.ACSEndpoint == "" {
			return fmt.Errorf("ACS_ENDPOINT is required when NOTIFIER_TYPE=acs")
		}
		if c.ACSSenderAddress == "" {
			return fmt.Errorf("ACS_SENDER_ADDRESS is required when NOTIFIER_TYPE=acs")
		}
		if c.FrontendBaseURL == "" {
			return fmt.Errorf("FRONTEND_BASE_URL is required when NOTIFIER_TYPE=acs (used for password-reset links)")
		}
		mode := strings.ToLower(c.ACSAuthMode)
		if mode != "managed_identity" && mode != "connection_string" {
			return fmt.Errorf("ACS_AUTH_MODE must be 'managed_identity' or 'connection_string'")
		}
		if mode == "connection_string" && c.ACSConnectionStringFile == "" && c.ACSConnectionStringInline == "" {
			return fmt.Errorf("ACS_CONNECTION_STRING_FILE or ACS_CONNECTION_STRING is required when ACS_AUTH_MODE=connection_string")
		}
	default:
		return fmt.Errorf("NOTIFIER_TYPE must be 'mock' or 'acs', got %q", c.NotifierType)
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
