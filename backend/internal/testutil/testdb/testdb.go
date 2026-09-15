//go:build integration

// Package testdb creates isolated, migrated PostgreSQL databases for integration tests.
package testdb

import (
	"context"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openschool-org/openschool/internal/database"
)

// Open creates a temporary database, applies every embedded migration, and
// registers cleanup. Tests are skipped when TEST_DATABASE_URL is not set.
func Open(t testing.TB) *pgxpool.Pool {
	t.Helper()

	rawURL := os.Getenv("TEST_DATABASE_URL")
	if rawURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Scheme != "postgres" {
		t.Fatalf("TEST_DATABASE_URL must be a valid postgres:// URL: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	adminURL := *parsed
	adminURL.Path = "/postgres"
	admin, err := pgx.Connect(ctx, adminURL.String())
	if err != nil {
		t.Fatalf("connect to PostgreSQL admin database: %v", err)
	}

	databaseName := "openschool_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	identifier := pgx.Identifier{databaseName}.Sanitize()
	if _, err := admin.Exec(ctx, "CREATE DATABASE "+identifier); err != nil {
		_ = admin.Close(ctx)
		t.Fatalf("create integration database: %v", err)
	}

	var pool *pgxpool.Pool
	t.Cleanup(func() {
		if pool != nil {
			pool.Close()
		}
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cleanupCancel()
		if _, err := admin.Exec(cleanupCtx, "DROP DATABASE IF EXISTS "+identifier+" WITH (FORCE)"); err != nil {
			t.Errorf("drop integration database: %v", err)
		}
		if err := admin.Close(cleanupCtx); err != nil {
			t.Errorf("close PostgreSQL admin connection: %v", err)
		}
	})

	testURL := *parsed
	testURL.Path = "/" + databaseName
	if err := database.RunMigrations(testURL.String()); err != nil {
		t.Fatalf("migrate integration database: %v", err)
	}
	pool, err = pgxpool.New(ctx, testURL.String())
	if err != nil {
		t.Fatalf("open integration database pool: %v", err)
	}
	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("ping integration database: %v", err)
	}
	return pool
}
