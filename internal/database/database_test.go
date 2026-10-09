package database

import (
	"context"
	"os"
	"testing"
	"time"
)

// TestConnectAndMigrate exercises the Postgres service. It is skipped unless
// TEST_DATABASE_URL points at a reachable database, so the unit test suite
// stays runnable without the service.
func TestConnectAndMigrate(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("set TEST_DATABASE_URL to run the Postgres integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := Connect(ctx, url)
	if err != nil {
		t.Fatalf("Connect: %v", err)
	}
	defer pool.Close()

	if err := Migrate(ctx, pool); err != nil {
		t.Fatalf("Migrate: %v", err)
	}

	var exists bool
	if err := pool.QueryRow(ctx,
		"SELECT to_regclass('public.content') IS NOT NULL",
	).Scan(&exists); err != nil {
		t.Fatalf("check content table: %v", err)
	}
	if !exists {
		t.Fatal("content table was not created")
	}
}
