package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	createTables()
}

func createTables() {
	query := `
    CREATE TABLE IF NOT EXISTS incidents (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        description TEXT,
		severity TEXT CHECK(severity IN ('Low', 'Medium', 'High')) NOT NULL,
		status TEXT CHECK(status IN ('Open', 'In Progress', 'Resolved')) NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		resolved_at DATETIME DEFAULT NULL
    );`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
