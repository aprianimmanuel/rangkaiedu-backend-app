package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/aprianimmanuel/rangkaiedu-backend/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// DB represents the database connection wrapper
type DB struct {
	Connection *sql.DB
	Config     *config.Config
}

// New creates a new database instance
func New(cfg *config.Config) (*DB, error) {
	return &DB{
		Config: cfg,
	}, nil
}

// Init initializes the database connection pool
func (db *DB) Init() error {
	// Add diagnostic logging
	log.Printf("[DEBUG] Database Init called with Config: %+v", db.Config)
	
	// Check if DSN method exists
	if db.Config == nil {
		return fmt.Errorf("config is nil")
	}
	
	// Check if DSN method exists by calling it
	dsn := db.Config.DSN()
	log.Printf("[DEBUG] DSN() returned: %s", dsn)
	
	if dsn == "" {
		return fmt.Errorf("DSN is empty")
	}
	
	connection, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool parameters
	connection.SetMaxOpenConns(25)
	connection.SetMaxIdleConns(25)
	connection.SetConnMaxLifetime(5 * 60 * time.Second)

	// Test the connection
	if err := connection.Ping(); err != nil {
		return err
	}

	db.Connection = connection
	log.Println("Database connection established successfully")
	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.Connection != nil {
		return db.Connection.Close()
	}
	return nil
}

// GetConnection returns the database connection
func (db *DB) GetConnection() *sql.DB {
	return db.Connection
}

// Global database instance
var globalDB *DB

// Init initializes the global database connection pool
func Init() error {
	cfg := config.Load()
	dbInstance, err := New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create database instance: %w", err)
	}
	
	if err := dbInstance.Init(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	
	globalDB = dbInstance
	log.Println("Global database connection pool initialized successfully")
	return nil
}

// Close closes the global database connection
func Close() error {
	if globalDB != nil {
		return globalDB.Close()
	}
	return nil
}

// GetDB returns the global database connection
func GetDB() *sql.DB {
	if globalDB == nil {
		return nil
	}
	return globalDB.GetConnection()
}