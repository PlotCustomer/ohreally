package config

import "testing"

func TestLoadRequiresDatabaseURL(t *testing.T) {
	t.Setenv("EPILEPTIC_DATABASE_URL", "")

	if _, err := Load(); err == nil {
		t.Fatal("expected an error when EPILEPTIC_DATABASE_URL is unset")
	}
}

func TestLoadAppliesDefaults(t *testing.T) {
	t.Setenv("EPILEPTIC_DATABASE_URL", "postgres://example")
	t.Setenv("EPILEPTIC_HTTP_ADDR", "")
	t.Setenv("EPILEPTIC_ENV", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.HTTPAddr != ":8080" {
		t.Errorf("HTTPAddr = %q, want :8080", cfg.HTTPAddr)
	}
	if cfg.Env != "debug" {
		t.Errorf("Env = %q, want debug", cfg.Env)
	}
}

func TestLoadRejectsUnknownEnv(t *testing.T) {
	t.Setenv("EPILEPTIC_DATABASE_URL", "postgres://example")
	t.Setenv("EPILEPTIC_ENV", "staging")

	if _, err := Load(); err == nil {
		t.Fatal("expected an error for an unknown EPILEPTIC_ENV")
	}
}
