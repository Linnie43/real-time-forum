package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes the SQLite database and returns a database connection object
func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Println("Failed to initialize the database")
		log.Fatal(err)
	}

	_, err = db.Exec(CreateTables)
	if err != nil {
		log.Println("Failed to create tables")
		log.Fatal(err)
	}

	return db
}
