package repository

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/amirdashtii/go_auth/config"
	"github.com/amirdashtii/go_auth/infrastructure/logger"
	"github.com/amirdashtii/go_auth/internal/core/errors"
	"github.com/amirdashtii/go_auth/internal/core/ports"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type PGRepository struct {
	db *sql.DB
}

func runMigrations(config *config.Config) {
	logfile, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	loggerConfig := ports.LoggerConfig{
		Level: "info",
		Environment: "development",
		ServiceName: "go_auth",
		Output: logfile,
	}
	logger := logger.NewZerologLogger(loggerConfig)


	m, err := migrate.New(
		"file://migrations",
		"postgres://"+config.DB.User+":"+config.DB.Password+"@"+config.DB.Host+":"+config.DB.Port+"/"+config.DB.Name+"?sslmode=disable",
	)
	if err != nil {
		logger.Fatal("Failed to create migrate instance", 
			ports.F("error", err),
		)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatal("Failed to run migrations", 
			ports.F("error", err),
		)
	}
}

func NewPGRepository() (*PGRepository, error) {
	logfile, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	loggerConfig := ports.LoggerConfig{
		Level: "info",
		Environment: "development",
		ServiceName: "go_auth",
		Output: logfile,
	}
	logger := logger.NewZerologLogger(loggerConfig)
	

	config, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	runMigrations(config)

	host := config.DB.Host
	user := config.DB.User
	password := config.DB.Password
	dbName := config.DB.Name
	port := config.DB.Port

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbName, port)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error("Failed to open database", 
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

	return &PGRepository{db: db}, nil
}

func (r *PGRepository) DB() *sql.DB {
	return r.db
}
