package main

import "database/sql"

type Database struct {
	db                  *sql.DB
	name                string
	numOfTables         int
	nameOfTables        []string
	tables              []Table
	totalSpace          int64
}

type Table struct {
	database *Database
	name     string
	columnsnames  []string
	columnTypes []string
	numOfRows     int64
	totalSpace    int64
}

type Model struct {
	database   *Database
	cursor     int              // which table we're pointing at
	selected   map[int]struct{} // which tables are expanded/selected
	viewMode   string           // "overview" or "table"
}


