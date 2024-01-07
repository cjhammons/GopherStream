package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func InitializeDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./goaudio.db")
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
			file_format TEXT NOT NULL,
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
		log.Printf("Executing query: %v", query)
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

func InsertArtist(db *sql.DB, name string, artFilePath string) error {
	var artistID int64

	// Check if artist already exists
	exists_query := `
		SELECT
			id
		FROM
			artists
		WHERE
			name = ?
	`
	err := db.QueryRow(exists_query, name).Scan(&artistID)
	if err == nil {
		log.Printf("Artist already exists: %v", name)
		return nil
	} else if err != sql.ErrNoRows {
		log.Printf("Error querying database: %v", err)
		return err
	}

	// Start a new transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare the statement
	stmt, err := tx.Prepare(`
		INSERT INTO 
			artists(name, art_file_path) 
			VALUES(?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the statement
	_, err = stmt.Exec(name, artFilePath)
	if err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func InsertSong(db *sql.DB, title string, artistID int64, albumID int64, genreID int64, trackNumber int64, filePath string, fileFormat string) error {
	var songId int64

	// Check if song already exists
	// We also need to check the artist id, as
	exists_query := `
	SELECT 
		id 
	FROM 
		artists 
	WHERE 
		title = ?
		AND 
		artist_id = ?
		AND
		album_id = ?
		AND
		genre_id = ?`
	err := db.QueryRow(exists_query, title, artistID, albumID, genreID).Scan(&songId)
	if err == nil {
		log.Printf("Song already exists: %v by %v", title, artistID)
		return nil
	} else if err != sql.ErrNoRows {
		log.Printf("Error querying database: %v", err)
		return err
	}

	// Start a new transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare the statement
	stmt, err := tx.Prepare(
		`INSERT INTO 
			songs(title, artist_id, album_id, genre_id, track_number, file_path, file_format) 
			VALUES(?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the statement
	_, err = stmt.Exec(title, artistID, albumID, genreID, trackNumber, filePath)
	if err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func InsertAlbum(db *sql.DB, title string, artistID int64, releaseDate string, artFilePath string) error {
	// Start a new transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare the statement
	stmt, err := tx.Prepare(`
		INSERT INTO 
			albums(title, artist_id, release_date, art_file_path) 
			VALUES(?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the statement
	_, err = stmt.Exec(title, artistID, releaseDate, artFilePath)
	if err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func InsertGenre(db *sql.DB, name string) error {
	// Start a new transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare the statement
	stmt, err := tx.Prepare(`
		INSERT INTO 
			genres(name) 
			VALUES(?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the statement
	_, err = stmt.Exec(name)
	if err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

//TODO: Add a function to insert a new record into the database
//TODO: Add a function to retrieve a record from the database
//TODO: Add a function to update a record in the database
