package database

import (
	"database/sql"
	"log"

	"github.com/dhowden/tag"
	_ "github.com/mattn/go-sqlite3"
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

func InsertOrUpdateArtist(db *sql.DB, metadata tag.Metadata) (int64, error) {
	var artistID int64

	// Check if artist already exists
	existsQuery := `SELECT id FROM artists WHERE name = ?`
	err := db.QueryRow(existsQuery, metadata.Artist()).Scan(&artistID)
	if err == nil {
		log.Printf("Artist already exists: %v", metadata.Artist())
		return artistID, nil
	} else if err != sql.ErrNoRows {
		return 0, err
	}

	// Insert new artist
	insertQuery := `
		INSERT INTO 
			artists(name, art_file_path)
		VALUES(?, ?)
	`
	result, err := db.Exec(insertQuery, metadata.Artist(), "")
	if err != nil {
		return 0, err
	}

	artistID, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return artistID, nil
}

func InsertOrUpdateAlbum(db *sql.DB, metadata tag.Metadata, artistID int64) (int64, error) {
	var albumID int64
	// Check if album already exists
	existsQuery := `
		SELECT 
			id 
		FROM albums 
		WHERE 
			title = ? 
			AND 
			artist_id = ?
	`
	err := db.QueryRow(existsQuery, metadata.Album(), artistID).Scan(&albumID)
	if err == nil {
		log.Printf("Album already exists: \"%v\"", metadata.Album())

		// Update existing album
		updateQuery := `
            UPDATE albums 
            SET 
				title = ?,
				artist_id = ?,
                release_date = ?
            WHERE 
                id = ?`
		_, err := db.Exec(updateQuery, metadata.Album(), artistID, metadata.Year(), albumID)
		if err != nil {
			return 0, err
		}
		return albumID, nil
	} else if err != sql.ErrNoRows {
		return 0, err
	}

	// Insert new album
	insertQuery := `
		INSERT INTO 
			albums (title, artist_id, release_date, art_file_path) 
			VALUES (?, ?, ?, ?)
	`
	result, err := db.Exec(insertQuery, metadata.Album(), artistID, metadata.Year(), "")
	if err != nil {
		return 0, err
	}

	albumID, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return albumID, nil
}

func InsertOrUpdateAlbumArt(db *sql.DB, metadata tag.Metadata, albumID int64) error {
	imageFilePath, err := SaveAlbumArt(metadata.Picture(), albumID)
	if err != nil {
		log.Printf("Error saving album art: %v", err)
	}

	query := `
		UPDATE albums
		SET
			art_file_path = ?
		WHERE
			id = ?
	`
	_, err = db.Exec(query, imageFilePath, albumID)
	if err != nil {
		return err
	}
	return nil
}

func InsertOrUpdateGenre(db *sql.DB, metadata tag.Metadata) (int64, error) {
	var genreID int64

	// Check if genre already exists
	existsQuery := `
		SELECT 
			id 
		FROM 
			genres 
		WHERE name = ?
	`
	err := db.QueryRow(existsQuery, metadata.Genre()).Scan(&genreID)
	if err == nil {
		log.Printf("Genre already exists: \"%v\"", metadata.Genre())
		return genreID, nil
	} else if err != sql.ErrNoRows {
		return -1, err
	}

	// Insert new genre
	insertQuery := `
		INSERT INTO 
			genres(name) 
		VALUES(?)
	`
	result, err := db.Exec(insertQuery, metadata.Genre())
	if err != nil {
		return -1, err
	}

	genreID, err = result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return genreID, nil
}

func InsertOrUpdateSong(db *sql.DB, metadata tag.Metadata, albumID int64, artistID int64, genreID int64, filePath string) (int64, error) {
	var songID int64
	trackNumber, _ := metadata.Track()

	// Check if song already exists
	existsQuery := `
		SELECT 
			id 
		FROM 
			songs 
		WHERE 
			title = ? 
			AND 
			album_id = ?
			AND
			artist_id = ?
	`
	err := db.QueryRow(existsQuery, metadata.Title(), albumID, artistID).Scan(&songID)
	if err == nil {
		log.Printf("Song \"%v\" already exists. Updating...", metadata.Title())

		// Update existing song
		updateQuery := `
            UPDATE songs 
            SET 
                artist_id = ?, 
                genre_id = ?, 
				album_id = ?,
                track_number = ?, 
                file_path = ?, 
                file_format = ? 
            WHERE 
                id = ?`
		_, err := db.Exec(updateQuery, artistID, genreID, albumID, trackNumber, filePath, metadata.FileType(), songID)
		if err != nil {
			return -1, err
		}

		return songID, nil
	} else if err != sql.ErrNoRows {
		return -1, err
	}

	// Insert new song
	insertQuery := `
		INSERT INTO 
			songs(title, album_id, artist_id, genre_id, track_number, file_path, file_format) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(insertQuery, metadata.Title(), albumID, artistID, genreID, trackNumber, filePath, metadata.FileType())
	if err != nil {
		return -1, err
	}

	songID, err = result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return songID, nil
}
