package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// DB wraps database connection
type DB struct {
	*sql.DB
}

// NewDB creates new database connection
func NewDB(dsn string) (*DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

// Close closes database connection
func (db *DB) Close() error {
	return db.DB.Close()
}
