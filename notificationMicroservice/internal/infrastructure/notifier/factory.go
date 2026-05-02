package notifier

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Config drives factory.New — it intentionally mirrors the env-var shape
// from infrastructure/config so wiring stays trivial.
type Config struct {
	Type                     string // "mock" or "acs"
	ACSEndpoint              string
	ACSSenderAddress         string
	ACSAuthMode              string // "managed_identity" or "connection_string"
	ACSConnectionStringFile  string // optional path containing the connection string
	ACSConnectionStringInline string // optional inline connection string (lower priority than file)
	FromName                 string
}

// New constructs a concrete Notifier from cfg. Mock and ACS are the
// supported types; anything else is an error.
func New(cfg Config) (Notifier, error) {
	switch strings.ToLower(cfg.Type) {
	case "", "mock":
		return NewMock(), nil
	case "acs":
		connStr, err := loadConnectionString(cfg)
		if err != nil {
			return nil, err
		}
		return NewACS(ACSConfig{
			Endpoint:         cfg.ACSEndpoint,
			Sender:           cfg.ACSSenderAddress,
			FromName:         cfg.FromName,
			AuthMode:         cfg.ACSAuthMode,
			ConnectionString: connStr,
		})
	default:
		return nil, fmt.Errorf("notifier: unknown type %q", cfg.Type)
	}
}

// loadConnectionString prefers a mounted secret file (Bicep / docker secret
// pattern) over an inline value. For managed-identity mode the connection
// string is unused, so we return early.
func loadConnectionString(cfg Config) (string, error) {
	if strings.ToLower(cfg.ACSAuthMode) != "connection_string" {
		return "", nil
	}
	if cfg.ACSConnectionStringFile != "" {
		raw, err := os.ReadFile(cfg.ACSConnectionStringFile)
		if err != nil {
			return "", fmt.Errorf("read ACS_CONNECTION_STRING_FILE: %w", err)
		}
		return strings.TrimSpace(string(raw)), nil
	}
	if cfg.ACSConnectionStringInline == "" {
		return "", errors.New("connection_string auth requires ACS_CONNECTION_STRING_FILE or ACS_CONNECTION_STRING")
	}
	return cfg.ACSConnectionStringInline, nil
}
