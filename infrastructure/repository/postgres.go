package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/amirdashtii/go_auth/config"
	"github.com/amirdashtii/go_auth/infrastructure/repository/migrations"
)

type PGRepository struct {
	db *sql.DB
}

func NewPGRepository() (*PGRepository, error) {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

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

	if err := migrations.RunMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	return &PGRepository{db: db}, nil
}
