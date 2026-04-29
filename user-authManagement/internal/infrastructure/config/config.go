// Package config loads runtime configuration from the environment.
// This is the only place in the codebase that reads ENV vars directly.
package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type HTTP struct {
	Port             int           `envconfig:"HTTP_PORT" default:"8080"`
	ReadTimeout      time.Duration `envconfig:"HTTP_READ_TIMEOUT" default:"10s"`
	WriteTimeout     time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" default:"15s"`
	ShutdownTimeout  time.Duration `envconfig:"HTTP_SHUTDOWN_TIMEOUT" default:"15s"`
	CORSAllowedOrigins []string    `envconfig:"CORS_ALLOWED_ORIGINS" default:""`
}

type Logging struct {
	Level  string `envconfig:"LOG_LEVEL"  default:"info"`
	Format string `envconfig:"LOG_FORMAT" default:"json"`
}

type Mongo struct {
	URI            string        `envconfig:"MONGO_URI" required:"true"`
	Database       string        `envconfig:"MONGO_DATABASE" default:"auth"`
	ConnectTimeout time.Duration `envconfig:"MONGO_CONNECT_TIMEOUT" default:"10s"`
}

type RabbitMQ struct {
	URL              string        `envconfig:"RABBITMQ_URL" required:"true"`
	Exchange         string        `envconfig:"RABBITMQ_EXCHANGE" default:"user.events"`
	ReconnectBackoff time.Duration `envconfig:"RABBITMQ_RECONNECT_BACKOFF" default:"5s"`
}

type JWT struct {
	// Inline PEM (single line, "\n" allowed as escape).
	PrivateKey string `envconfig:"JWT_PRIVATE_KEY"`
	PublicKey  string `envconfig:"JWT_PUBLIC_KEY"`
	// Or: path to a PEM file (preferred — avoids ENV escaping issues).
	PrivateKeyFile string        `envconfig:"JWT_PRIVATE_KEY_FILE"`
	PublicKeyFile  string        `envconfig:"JWT_PUBLIC_KEY_FILE"`
	Issuer         string        `envconfig:"JWT_ISSUER" default:"user-auth"`
	AccessTTL      time.Duration `envconfig:"JWT_ACCESS_TTL"  default:"15m"`
	RefreshTTL     time.Duration `envconfig:"JWT_REFRESH_TTL" default:"336h"`
}

type Auth struct {
	PasswordMinLength       int           `envconfig:"PASSWORD_MIN_LENGTH" default:"12"`
	BcryptCost              int           `envconfig:"BCRYPT_COST" default:"12"`
	LoginRateLimitPerMin    int           `envconfig:"LOGIN_RATE_LIMIT_PER_MIN" default:"5"`
	RegisterRateLimitPerMin int           `envconfig:"REGISTER_RATE_LIMIT_PER_MIN" default:"5"`
	PasswordResetTTL        time.Duration `envconfig:"PASSWORD_RESET_TTL" default:"30m"`

	// Optional bootstrap admin. Both must be set together; the seed runs
	// idempotently on startup. Password must satisfy PasswordMinLength.
	AdminBootstrapEmail    string `envconfig:"ADMIN_BOOTSTRAP_EMAIL"`
	AdminBootstrapPassword string `envconfig:"ADMIN_BOOTSTRAP_PASSWORD"`
}

type Config struct {
	HTTP     HTTP
	Logging  Logging
	Mongo    Mongo
	RabbitMQ RabbitMQ
	JWT      JWT
	Auth     Auth
}

// Load reads ENV vars (and a .env file if present in the working dir,
// loaded by the caller) and returns a populated Config.
func Load() (*Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	if err := c.JWT.resolveKeys(); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	if (c.Auth.AdminBootstrapEmail == "") != (c.Auth.AdminBootstrapPassword == "") {
		return nil, errors.New("config: ADMIN_BOOTSTRAP_EMAIL and ADMIN_BOOTSTRAP_PASSWORD must be set together")
	}
	return &c, nil
}

// resolveKeys loads PEM material from the configured source.
// Prefers *_FILE over inline ENV — files avoid all newline-escaping issues.
func (j *JWT) resolveKeys() error {
	if j.PrivateKeyFile != "" {
		b, err := os.ReadFile(j.PrivateKeyFile)
		if err != nil {
			return fmt.Errorf("read JWT_PRIVATE_KEY_FILE: %w", err)
		}
		j.PrivateKey = string(b)
	}
	if j.PublicKeyFile != "" {
		b, err := os.ReadFile(j.PublicKeyFile)
		if err != nil {
			return fmt.Errorf("read JWT_PUBLIC_KEY_FILE: %w", err)
		}
		j.PublicKey = string(b)
	}
	if j.PrivateKey == "" || j.PublicKey == "" {
		return errors.New("JWT keys missing: set JWT_PRIVATE_KEY/PUBLIC_KEY or JWT_PRIVATE_KEY_FILE/PUBLIC_KEY_FILE")
	}
	return nil
}
