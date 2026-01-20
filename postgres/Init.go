package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func InitDB(pgConnStr string) (*sql.DB, error) {
	var db *sql.DB
	var err error
	// Retry Postgres connection 5 times
	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", pgConnStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("PostgreSQL connection successful", "URL", maskDatabaseURL(pgConnStr), "attempt", i+1)
				break
			}
		}
		log.Println("Waiting for PostgreSQL...", "attempt", i+1, "max_attempts", 5, "error", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Println("Failed to connect to PostgreSQL", "URL", maskDatabaseURL(pgConnStr), "error", err)
		return nil, fmt.Errorf("Failed to connect to PostgreSQL: %w", err)
	}

	// Test PostgreSQL connection
	if err := db.Ping(); err != nil {
		log.Println("Failed to ping PostgreSQL", "URL", maskDatabaseURL(pgConnStr), "error", err)
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	log.Println("PostgreSQL connection pool configured", "URL", maskDatabaseURL(pgConnStr), "max_open_conns", 25, "max_idle_conns", 5, "conn_max_lifetime", "5m")
	return db, nil
}
