package main

import "database/sql"

type Database struct {
	db                  *sql.DB
	name                string
	numOfTables         int
	nameOfTables        []string
	totalSpace          int
}

type Table struct {
	database *Database
	name     string
	columnsnames  []string
	columnTypes []string
	numOfRows     int
	totalSpace    int
}
