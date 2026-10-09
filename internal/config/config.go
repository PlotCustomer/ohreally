// Package config loads the runtime configuration for the Epileptic
// application component from the environment.
package config

import (
	"fmt"
	"os"
	"strings"
)

// Default values for the local Postgres service declared in
// .plot/package.json.
const (
	defaultPort     = "8080"
	defaultHost     = "localhost"
	defaultPortDB   = "5432"
	defaultUser     = "epileptic"
	defaultPassword = "epileptic"
	defaultDatabase = "epileptic"
	defaultSSLMode  = "disable"
	databaseURLEnv  = "DATABASE_URL"
	portEnv         = "PORT"
	hostEnv         = "PGHOST"
	portDBEnv       = "PGPORT"
	userEnv         = "PGUSER"
	passwordEnv     = "PGPASSWORD"
	databaseNameEnv = "PGDATABASE"
	sslModeEnv      = "PGSSLMODE"
)

// Config is the runtime configuration of the Epileptic service.
type Config struct {
	// Port is the TCP port the Communication Service listens on.
	Port string
	// DatabaseURL is the connection string of the Postgres service.
	DatabaseURL string
}

// Load reads the configuration from the environment, falling back to the
// defaults of the local Postgres service.
func Load() Config {
	return Config{
		Port:        getenv(portEnv, defaultPort),
		DatabaseURL: databaseURL(),
	}
}

// databaseURL returns an explicit DATABASE_URL when set, otherwise it
// assembles one from the individual Postgres environment variables.
func databaseURL() string {
	if url := strings.TrimSpace(os.Getenv(databaseURLEnv)); url != "" {
		return url
	}

	host := getenv(hostEnv, defaultHost)
	port := getenv(portDBEnv, defaultPortDB)
	user := getenv(userEnv, defaultUser)
	password := getenv(passwordEnv, defaultPassword)
	name := getenv(databaseNameEnv, defaultDatabase)
	sslMode := getenv(sslModeEnv, defaultSSLMode)

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, name, sslMode,
	)
}

func getenv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
