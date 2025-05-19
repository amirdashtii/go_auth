package repository

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	"github.com/amirdashtii/go_auth/config"
	"github.com/amirdashtii/go_auth/infrastructure/logger"
	"github.com/amirdashtii/go_auth/internal/core/errors"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type PGRepository struct {
	db     *sql.DB
	logger ports.Logger
}

var (
	pgOnce          sync.Once
	pgRepository  *PGRepository
)

func GetPGRepository(config *config.Config) (*PGRepository, error) {
	var err error
	pgOnce.Do(func() {
		pgRepository, err = newPGRepository(config)
	})
	return pgRepository, err
}

func newPGRepository(config *config.Config) (*PGRepository, error) {
	loggerConfig := ports.LoggerConfig{
		Level:       "info",
		Environment: config.Environment,
		ServiceName: "go_auth",
		Output:      os.Stdout,
	}
	logger := logger.NewZerologLogger(loggerConfig)

	host := config.DB.Host
	user := config.DB.User
	password := config.DB.Password
	dbName := config.DB.Name
	port := config.DB.Port

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbName, port)
	
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("Failed to open database connection",
			ports.F("error", err),
		)
		return nil, errors.ErrDatabaseInit
	}

	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping database",
			ports.F("error", err),
		)
		return nil, errors.ErrDatabaseInit
	}

	return &PGRepository{
		db:     db,
		logger: logger,
	}, nil
}

// DB returns the underlying database connection
func (r *PGRepository) DB() *sql.DB {
	return r.db
}

// Close closes the database connection
func (r *PGRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}
