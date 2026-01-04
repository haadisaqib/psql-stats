package main

import (
	"database/sql"
	"fmt"
	"net/url"
	_ "github.com/lib/pq"
	"sync"
)

func fetchDatabaseInfo(db *Database, connectionString string) error {
	var wg sync.WaitGroup
	var errs = make(chan error, 4)
	wg.Add(4)

	go func() {
		db.name = GetDatabaseName(connectionString, &wg)
	}()

	go func() {
		if err := GetNumberOfTables(db, &wg); err != nil {
			errs <- err
		}
	}()

	go func() {
		if err := GetTableNames(db, &wg); err != nil {
			errs <- err
		}
	}()
	go func() {
		if err := TotalSpace(db, &wg); err != nil {
			errs <- err
		}
	}()

	wg.Wait()
	close(errs)

	return nil
}

func GetDatabaseName(connectionString string, wg *sync.WaitGroup) string {
	defer wg.Done()
	u, err := url.Parse(connectionString)
	if err != nil {
		return ""
	}
	return u.Path[1:]
}

func GetNumberOfTables(db *Database, wg *sync.WaitGroup) error {
	defer wg.Done()
	query := "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'"
	rows, err := db.db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to get number of tables: %v", err)
	}
	defer rows.Close()
	if rows.Next() {
		if err := rows.Scan(&db.numOfTables); err != nil {
			return err
		}
	}
	return rows.Err()
}


func GetTableNames(db *Database, wg *sync.WaitGroup) error {
	defer wg.Done()
	query := "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'"
	rows, err := db.db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to get table names: %v", err)
	}
	defer rows.Close()
	db.nameOfTables = []string{}
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return err
		}
		db.nameOfTables = append(db.nameOfTables, tableName)
	}
	return rows.Err()
}

func TotalSpace(db *Database, wg *sync.WaitGroup) error {
	defer wg.Done()
	query := `
		SELECT SUM(pg_total_relation_size(pg_class.oid))
		FROM pg_class
		JOIN pg_namespace ON pg_namespace.oid = pg_class.relnamespace
		WHERE pg_class.relkind = 'r' AND pg_namespace.nspname = 'public'
	`
	rows, err := db.db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to get total space: %v", err)
	}
	defer rows.Close()

	if rows.Next() {
		var totalSpace sql.NullInt64
		if err := rows.Scan(&totalSpace); err != nil {
			return err
		}
		if totalSpace.Valid {
			db.totalSpace = totalSpace.Int64
		} else {
			db.totalSpace = 0
		}
	}
	return nil
}

func fetchTableInfo(table *Table) error {
	if err := GetColumnNames(table); err != nil {
		return err
	}
	if err := GetColumnTypes(table); err != nil {
		return err
	}
	if err := GetNumberOfRows(table); err != nil {
		return err
	}
	if err := GetTotalSpace(table); err != nil {
		return err
	}
	table.database.tables = append(table.database.tables, *table)
	return nil
}

func GetColumnNames(table *Table) error {
	query := "SELECT column_name FROM information_schema.columns WHERE table_name = $1"
	rows, err := table.database.db.Query(query, table.name)
	if err != nil {
		return fmt.Errorf("failed to get column names: %v", err)
	}
	defer rows.Close()
	table.columnsnames = []string{}
	for rows.Next() {
		var columnName string
		if err := rows.Scan(&columnName); err != nil {
			return err
		}
		table.columnsnames = append(table.columnsnames, columnName)
	}
	return rows.Err()
}

func GetColumnTypes(table *Table) error {
	query := "SELECT data_type FROM information_schema.columns WHERE table_name = $1"
	rows, err := table.database.db.Query(query, table.name)
	if err != nil {
		return fmt.Errorf("failed to get column types: %v", err)
	}
	defer rows.Close()
	table.columnTypes = []string{}
	for rows.Next() {
		var columnType string
		if err := rows.Scan(&columnType); err != nil {
			return err
		}
		table.columnTypes = append(table.columnTypes, columnType)
	}
	return rows.Err()
}

func GetNumberOfRows(table *Table) error {
	query := "SELECT COUNT(*) FROM " + table.name
	rows, err := table.database.db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to get number of rows: %v", err)
	}
	defer rows.Close()
	if rows.Next() {
		if err := rows.Scan(&table.numOfRows); err != nil {
			return err
		}
	}
	return rows.Err()
}

func GetTotalSpace(table *Table) error {
	query := "SELECT pg_total_relation_size('" + table.name + "')"
	rows, err := table.database.db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to get total space: %v", err)
	}
	defer rows.Close()
	if rows.Next() {
		if err := rows.Scan(&table.totalSpace); err != nil {
			return err
		}
	}
	return rows.Err()
}