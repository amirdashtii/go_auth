package repository

import (
	"database/sql"
	"fmt"
	"os"
)

type PGRepository struct {
	db *sql.DB
}

func NewPGRepository() (*PGRepository, error) {

	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DATABASE")
	port := os.Getenv("POSTGRES_PORT")

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
