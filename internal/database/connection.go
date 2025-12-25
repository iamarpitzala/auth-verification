package database

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection with a timeout
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
