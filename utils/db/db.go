package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/aprianimmanuel/rangkaiedu-backend/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

// Init initializes the database connection pool
func Init() error {
	cfg := config.Load()
	
	db, err := sql.Open("pgx", cfg.DSN())
	if err != nil {
		return err
	}

	// Set connection pool parameters
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * 60 * time.Second)

	// Test the connection
	if err := db.Ping(); err != nil {
		return err
	}

	DB = db
	log.Println("Database connection established successfully")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// GetDB returns the global database connection
func GetDB() *sql.DB {
	return DB
}