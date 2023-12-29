package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func InitializeDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Start a new transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to start transaction: %v", err)
	}

	create_table_queries := []string{
		`
		CREATE TABLE IF NOT EXISTS artists (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			art_file_path TEXT
		);
	`, `
		CREATE TABLE IF NOT EXISTS albums (
			id INTEGER PRIMARY KEY,
			title TEXT NOT NULL,
			artist_id INTEGER,
			release_date TEXT NOT NULL,
			art_file_path TEXT NOT NULL,
			FOREIGN KEY (artist_id) REFERENCES artists(id)
		);
	`, `
		CREATE TABLE IF NOT EXISTS songs (
			id INTEGER PRIMARY KEY,
			title TEXT NOT NULL,
			artist_id INTEGER NOT NULL,
			album_id INTEGER NOT NULL,
			genre_id INTEGER NOT NULL,
			track_number INTEGER NOT NULL,
			file_path TEXT NOT NULL,
			FOREIGN KEY (album_id) REFERENCES albums(id),
			FOREIGN KEY (genre_id) REFERENCES genre(id),
			FOREIGN KEY (artist_id) REFERENCES artists(id)
		);
	`, `
		CREATE TABLE IF NOT EXISTS genres (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL
		);
	`,
	}

	for _, query := range create_table_queries {
		if _, err := tx.Exec(query); err != nil {
			// If an error occurs, roll back the transaction and report the error
			tx.Rollback()
			log.Fatalf("Failed to create table: %v", err)
		}
	}
	// If everything is successful, commit the transaction
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	return db, nil
}

//TODO: Add a function to insert a new record into the database
//TODO: Add a function to retrieve a record from the database
//TODO: Add a function to update a record in the database
