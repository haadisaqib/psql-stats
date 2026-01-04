# psql-stats

A terminal UI tool for viewing PostgreSQL database statistics.

## Installation

Clone the repository and build:

```bash
git clone https://github.com/haadisaqib/psql-stats.git
cd psql-stats
go build .
```

## Usage

Run the program with a PostgreSQL connection string:

```bash
./psql-stats -connection "postgresql://user:password@host:port/database?sslmode=require"
```

## Navigation

- Up/Down arrows or k/j: Navigate through tables
- Enter or Space: Expand/collapse table details
- q: Quit

## Features

- View database name, number of tables, and total space usage
- Browse all tables with row counts and space usage
- Expand tables to see column names and types
- Clean, interactive terminal interface

## Requirements

- Go 1.20 or later
- PostgreSQL database access

## License

MIT

