package utils

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

// DB is a global variable for database connection pooling
var DB *sql.DB

// InitializeDatabase sets up the database connection
func InitializeDatabase() {
	dsn := os.Getenv("DB_DATA")
	var err error

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open MySQL connection: %v", err)
	}

	// Verify the connection with Ping
	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Successfully connected to the database.")
}
