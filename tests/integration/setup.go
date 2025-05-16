package integration

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/amirdashtii/go_auth/infrastructure/repository"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	testcontainers "github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	container        *tcpostgres.PostgresContainer
	db               *sql.DB
	authRepo         ports.AuthRepository
	containerStarted bool
	logger           ports.Logger
)

func cleanupDatabase() error {
	_, err := db.Exec("TRUNCATE TABLE users CASCADE")
	if err != nil {
		return fmt.Errorf("failed to truncate users table: %v", err)
	}
	return nil
}

func setupTestEnvironment() error {
	// Initialize logger
	logger = NewTestLogger("test")

	// Create PostgreSQL container
	ctx := context.Background()
	var err error
	container, err = tcpostgres.Run(ctx,
		"postgres:15-alpine",
		tcpostgres.WithDatabase("go_auth_test"),
		tcpostgres.WithUsername("go_auth"),
		tcpostgres.WithPassword("go_auth"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections")),
	)
	if err != nil {
		return fmt.Errorf("failed to start container: %v", err)
	}

	// Get container connection details
	host, err := container.Host(ctx)
	if err != nil {
		return fmt.Errorf("failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return fmt.Errorf("failed to get container port: %v", err)
	}

	// Set environment variables
	os.Setenv("DB_HOST", host)
	os.Setenv("DB_PORT", port.Port())
	os.Setenv("DB_USER", "go_auth")
	os.Setenv("DB_PASSWORD", "go_auth")
	os.Setenv("DB_NAME", "go_auth_test")

	// Initialize database connection
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port.Port(),
		"go_auth",
		"go_auth",
		"go_auth_test",
	)
	fmt.Printf("[DEBUG] DB Connection: %s\n", dsn)
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %v", err)
	}

	// Wait for 5 seconds to ensure database is ready
	time.Sleep(5 * time.Second)

	// Check if database is ready by executing a simple query
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}
	fmt.Println("[DEBUG] Database is ready.")

	// Run migrations
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %v", err)
	}
	migrationsPath := filepath.Join(wd, "..", "..", "migrations")
	fmt.Printf("[DEBUG] Migrations Path: %s\n", migrationsPath)

	m, err := migrate.New(
		"file://"+migrationsPath,
		"postgres://"+
			"go_auth"+":"+
			"go_auth"+"@"+
			host+":"+
			port.Port()+
			"/"+
			"go_auth_test"+"?sslmode=disable",
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	// Initialize repository
	authRepo = repository.NewPGAuthRepository(db, logger)
	containerStarted = true
	return nil
}

func cleanup() {
	if db != nil {
		db.Close()
	}
	if containerStarted {
		container.Terminate(context.Background())
	}
}

func init() {
	err := setupTestEnvironment()
	if err != nil {
		fmt.Printf("Failed to setup test environment: %v\n", err)
		os.Exit(1)
	}
}

func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()

	// Cleanup
	cleanup()

	os.Exit(code)
}
