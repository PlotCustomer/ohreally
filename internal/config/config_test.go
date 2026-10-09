package config

import (
	"strings"
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("DATABASE_URL", "")
	t.Setenv("PGHOST", "")
	t.Setenv("PGPORT", "")
	t.Setenv("PGUSER", "")
	t.Setenv("PGPASSWORD", "")
	t.Setenv("PGDATABASE", "")
	t.Setenv("PGSSLMODE", "")

	cfg := Load()

	if cfg.Port != defaultPort {
		t.Fatalf("Port = %q, want %q", cfg.Port, defaultPort)
	}
	want := "postgres://epileptic:epileptic@localhost:5432/epileptic?sslmode=disable"
	if cfg.DatabaseURL != want {
		t.Fatalf("DatabaseURL = %q, want %q", cfg.DatabaseURL, want)
	}
}

func TestLoadDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://custom:secret@db:5432/custom?sslmode=require")

	cfg := Load()

	if !strings.Contains(cfg.DatabaseURL, "custom:secret@db:5432/custom") {
		t.Fatalf("DatabaseURL = %q, want explicit DATABASE_URL honoured", cfg.DatabaseURL)
	}
}

func TestLoadPostgresParts(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("PGHOST", "db.internal")
	t.Setenv("PGPORT", "6543")
	t.Setenv("PGUSER", "app")
	t.Setenv("PGPASSWORD", "pw")
	t.Setenv("PGDATABASE", "content")
	t.Setenv("PGSSLMODE", "require")

	cfg := Load()

	want := "postgres://app:pw@db.internal:6543/content?sslmode=require"
	if cfg.DatabaseURL != want {
		t.Fatalf("DatabaseURL = %q, want %q", cfg.DatabaseURL, want)
	}
}
