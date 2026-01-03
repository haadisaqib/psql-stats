package main

import (
	"database/sql"
	"fmt"
	"log"
	"flag"
	"os"
	"strings"
	
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/lib/pq"
)

func ConnecttoDatabase(db *Database, connStr string) error {
	if connStr == "" {
		return fmt.Errorf("connection string is required")
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

// Init is called when the program starts
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.database.tables)-1 {
				m.cursor++
			}
		case "enter", " ":
			// Toggle selection/expansion
			if _, ok := m.selected[m.cursor]; ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}
	return m, nil
}

// View renders the UI
func (m Model) View() string {
	var s strings.Builder

	const pageSize = 10 
	numTables := len(m.database.tables)

	cursor := m.cursor
	if numTables > 0 && cursor >= numTables {
		cursor = numTables - 1
	} else if cursor < 0 {
		cursor = 0
	}

	pageStart := 0
	pageEnd := 0
	if numTables <= pageSize {
		pageStart = 0
		pageEnd = numTables
	} else {
		pageStart = cursor - pageSize/2
		if pageStart < 0 {
			pageStart = 0
		}
		pageEnd = pageStart + pageSize
		if pageEnd > numTables {
			pageEnd = numTables
			pageStart = pageEnd - pageSize
			if pageStart < 0 {
				pageStart = 0
			}
		}
	}

	// Database Overview Header
	s.WriteString("╔═══════════════════════════════════════════════════════════╗\n")
	s.WriteString("║                    DATABASE STATISTICS                    ║\n")
	s.WriteString("╚═══════════════════════════════════════════════════════════╝\n\n")

	// Database Info
	s.WriteString(fmt.Sprintf("Database Name:     %s\n", m.database.name))
	s.WriteString(fmt.Sprintf("Number of Tables:  %d\n", m.database.numOfTables))
	s.WriteString(fmt.Sprintf("Total Space:       %s\n\n", formatBytes(m.database.totalSpace)))

	// Tables Section
	s.WriteString("╔═══════════════════════════════════════════════════════════╗\n")
	s.WriteString("║                         TABLES                            ║\n")
	s.WriteString("╚═══════════════════════════════════════════════════════════╝\n\n")

	if numTables == 0 {
		s.WriteString("No tables found.\n")
	} else {
		for i := pageStart; i < pageEnd; i++ {
			table := m.database.tables[i]
			// Cursor indicator
			cursorStr := " "
			if m.cursor == i {
				cursorStr = ">"
			}

			// Expanded indicator
			expanded := " "
			if _, ok := m.selected[i]; ok {
				expanded = "▼"
			} else {
				expanded = "▶"
			}

			// Table summary line
			s.WriteString(fmt.Sprintf("%s %s %s\n", cursorStr, expanded, table.name))
			s.WriteString(fmt.Sprintf("   Rows: %d | Space: %s\n", table.numOfRows, formatBytes(table.totalSpace)))

			// Expanded details
			if _, ok := m.selected[i]; ok {
				s.WriteString("   ───────────────────────────────────────────────────────\n")
				s.WriteString("   Columns:\n")
				for j, colName := range table.columnsnames {
					colType := ""
					if j < len(table.columnTypes) {
						colType = table.columnTypes[j]
					}
					s.WriteString(fmt.Sprintf("     • %s (%s)\n", colName, colType))
				}
				s.WriteString("\n")
			}
		}
		// Show an indicator if not all tables are visible
		if pageStart > 0 {
			s.WriteString("(↑ More above...)\n")
		}
		if pageEnd < numTables {
			s.WriteString("(↓ More below...)\n")
		}
	}

	// Footer
	s.WriteString("\n")
	s.WriteString("─────────────────────────────────────────────────────────────\n")
	s.WriteString("↑/↓: Navigate  Enter/Space: Expand/Collapse  q: Quit\n")

	return s.String()
}

// formatBytes converts bytes to human-readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func main() {
	connectionString := flag.String("connection", "", "The connection string to the database")
	flag.Parse()

	if *connectionString == "" {
		log.Fatalf("PostgreSQL connection string is required")
	}

	db := &Database{
		db: nil,
		name: "",
		numOfTables: 0,
		nameOfTables: []string{},
		tables: []Table{},
		totalSpace: 0,
	}

	if err := ConnecttoDatabase(db, *connectionString); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := fetchDatabaseInfo(db, *connectionString); err != nil {
		log.Fatalf("Failed to fetch database info: %v", err)
	}

	for _, tableName := range db.nameOfTables {
		table := &Table{
			database: db,
			name: tableName,
		}
		if err := fetchTableInfo(table); err != nil {
			log.Fatalf("Failed to fetch table info for %s: %v", tableName, err)
		}
	}


	// Initialize the model
	initialModel := Model{
		database: db,
		cursor:   0,
		selected: make(map[int]struct{}),
		viewMode: "overview",
	}


	// Clear the terminal
	fmt.Print("\033[H\033[2J")

	// Start the Bubble Tea program
	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}