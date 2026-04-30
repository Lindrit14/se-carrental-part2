package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	RabbitMQURL  string
	HTTPPort     string
	NotifierType string // "mock" or "smtp"
}

func Load() (*Config, error) {
	c := &Config{
		RabbitMQURL:  getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		HTTPPort:     getEnv("HTTP_PORT", "8084"),
		NotifierType: getEnv("NOTIFIER_TYPE", "mock"),
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
	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
