package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)


func GetDatabaseName(connectionString string) string {
	re := regexp.MustCompile(`databaseName=([^;]+)`)
	match := re.FindStringSubmatch(connectionString)
	if len(match) >= 2 {
		return match[1]
	}
	return ""
}