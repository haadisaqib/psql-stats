package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)


/**
look at environment variables for the database connection
*/

func GetDatabaseConnectionString() string {
	return os.Getenv("DATABASE_CONNECTION_STRING")
}

func ConnecttoDatabase(db *Database) error {
	connStr := GetDatabaseConnectionString()
	if connStr == "" {
		return fmt.Errorf("DATABASE_CONNECTION_STRING environment variable is not set")
	}

	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	if err := sqlDB.Ping(); err != nil {
		return err
	}

	db.db = sqlDB
	return nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db := &Database{}
	err := ConnecttoDatabase(db)

	fmt.Println("Connected to database")

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}