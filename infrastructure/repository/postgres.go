package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/amirdashtii/go_auth/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type PGRepository struct {
	db *sql.DB
}

func runMigrations(config *config.Config) {
	m, err := migrate.New(
		"file://migrations",
		"postgres://"+config.DB.User+":"+config.DB.Password+"@"+config.DB.Host+":"+config.DB.Port+"/"+config.DB.Name+"?sslmode=disable",
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}
}

func NewPGRepository() (*PGRepository, error) {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
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
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PGRepository{db: db}, nil
}

func (r *PGRepository) DB() *sql.DB {
	return r.db
}
