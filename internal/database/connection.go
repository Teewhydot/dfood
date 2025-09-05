package database

import (
	"database/sql"
	"dfood/internal/config"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDatabase(cfg *config.Config) error {
	var err error
	DB, err = sql.Open(cfg.DB.Driver, cfg.DB.Datasource)
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("could not ping database: %w", err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	if err = createTables(); err != nil {
		return fmt.Errorf("could not create tables: %w", err)
	}

	return nil
}

func createTables() error {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		first_name TEXT NOT NULL,
		last_name TEXT,
		email TEXT NOT NULL,
		password TEXT NOT NULL
	);
	`
	
	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		start_time DATETIME NOT NULL,
		end_time DATETIME NOT NULL,
		user_id INTEGER NOT NULL
	);
	`

	if _, err := DB.Exec(createUsersTable); err != nil {
		return fmt.Errorf("could not create users table: %w", err)
	}

	if _, err := DB.Exec(createEventsTable); err != nil {
		return fmt.Errorf("could not create events table: %w", err)
	}

	return nil
}

func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}