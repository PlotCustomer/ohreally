// Package config holds the runtime configuration of the Epileptic
// application component.
package config

import (
	"fmt"
	"os"
	"strings"
)

// Config is the runtime configuration of the Epileptic application
// component, loaded from the environment.
type Config struct {
	// Env is the Gin run mode (debug, release or test).
	Env string
	// HTTPAddr is the address the Communication Service listens on.
	HTTPAddr string
	// DatabaseURL is the Postgres connection string used by the Content
	// Management function to access the Content data object.
	DatabaseURL string
}

// Load reads the configuration from the process environment, applying
// defaults suitable for local development.
func Load() (Config, error) {
	cfg := Config{
		Env:         valueOrDefault("EPILEPTIC_ENV", "debug"),
		HTTPAddr:    valueOrDefault("EPILEPTIC_HTTP_ADDR", ":8080"),
		DatabaseURL: strings.TrimSpace(os.Getenv("EPILEPTIC_DATABASE_URL")),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("config: EPILEPTIC_DATABASE_URL is required")
	}

	switch cfg.Env {
	case "debug", "release", "test":
	default:
		return Config{}, fmt.Errorf("config: EPILEPTIC_ENV must be one of debug, release or test, got %q", cfg.Env)
	}

	return cfg, nil
}

func valueOrDefault(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}
