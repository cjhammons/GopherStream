package storage

import (
	"database/sql"
	"dhawden/tag"
	"log"
	"os"
	"path/filepath"
)

var allowedFileFormats = []string{".mp3", ".flac"}

func ScanLibrary(directory string, db *sql.DB) error {
	log.Println("Beginning library scan...")
	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && isAllowedAudioFile(path) {
			if err := processAudioFile(path, db); err != nil {
				log.Printf("Error processing audio file %s: %v", path, err)
				// Optionally, you can choose to continue on error
				// return nil
			}
		}
		return nil
	})
}

func isAllowedAudioFile(path string) bool {
	extension := filepath.Ext(path)
	for _, allowedExtension := range allowedFileFormats {
		if extension == allowedExtension {
			return true
		}
	}
	return false
}

func processAudioFile(filePath string, db *sql.DB) error {
	log.Printf("Processing file: %v", filePath)
	metadata, err := extractMetadata(filePath)
	if err != nil {
		log.Printf("Error extracting metadata: %v", err)
		return err
	}

	// Insert metadata into database
	artistID, err := InsertOrUpdateArtist(db, metadata, filePath)
	if err != nil {
		return err
	}

	albumID, err := InsertOrUpdateAlbum(db, metadata, artistID, filePath)
	if err != nil {
		return err
	}

	// Insert or update genre and song
	genreID, err := InsertOrUpdateGenre(db, metadata)
	if err != nil {
		return err
	}

	songID, err := InsertOrUpdateSong(db, metadata, albumID, artistID, genreID, filePath)
	if err != nil {
		return err
	}

	return nil
}

func InsertOrUpdateArtist(db *sql.DB, metadata *tag.Metadata, filePath string) (int64, error) {
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
	insertQuery := `INSERT INTO artists(name, art_file_path) VALUES(?, ?)`
	result, err := db.Exec(insertQuery, metadata.Artist(), filePath)
	if err != nil {
		return 0, err
	}

	artistID, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return artistID, nil
}

func InsertOrUpdateAlbum(db *sql.DB, metadata *tag.Metadata, artistID int64, filePath string) (int64, error) {
	var albumID int64

	// Check if album already exists
	existsQuery := `SELECT id FROM albums WHERE title = ? AND artist_id = ?`
	err := db.QueryRow(existsQuery, metadata.Album(), artistID).Scan(&albumID)
	if err == nil {
		log.Printf("Album already exists: %v", metadata.Album())
		return albumID, nil
	} else if err != sql.ErrNoRows {
		return 0, err
	}

	// Insert new album
	insertQuery := `INSERT INTO albums (title, artist_id, release_date, art_file_path) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(insertQuery, metadata.Album(), artistID, metadata.Year(), filePath)
	if err != nil {
		return 0, err
	}

	albumID, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return albumID, nil
}

func InsertOrUpdateGenre(db *sql.DB, metadata *tag.Metadata) error {
	var genreID int64

	// Check if genre already exists
	existsQuery := `SELECT id FROM genres WHERE name = ?`
	err := db.QueryRow(existsQuery, metadata.Genre()).Scan(&genreID)
	if err == nil {
		log.Printf("Genre already exists: %v", metadata.Genre())
		return nil
	} else if err != sql.ErrNoRows {
		return err
	}

	// Insert new genre
	insertQuery := `INSERT INTO genres(name) VALUES(?)`
	result, err := db.Exec(insertQuery, metadata.Genre())
	if err != nil {
		return err
	}

	genreID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func InsertOrUpdateSong(db *sql.DB, metadata *tag.Metadata, albumID, artistID int64, filePath string) error {
	var songID int64
	trackNumber, trackTotal := metadata.Track()

	// Check if song already exists
	existsQuery := `SELECT id FROM songs WHERE title = ? AND album_id = ?`
	err := db.QueryRow(existsQuery, metadata.Title(), albumID).Scan(&songID)
	if err == nil {
		log.Printf("Song already exists: %v", metadata.Title())
		return nil
	} else if err != sql.ErrNoRows {
		return err
	}

	// Insert new song
	insertQuery := `INSERT INTO songs(title, album_id, artist_id, track_number, track_total, file_path) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(insertQuery, metadata.Title(), albumID, artistID, trackNumber, trackTotal, filePath)
	if err != nil {
		return err
	}

	songID, err = result.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

func extractMetadata(filePath string) (*tag.Metadata, error) {
	m, err := tag.ReadFrom(filePath)
	if err != nil {
		return nil, err
	}
	log.Print(m.Format()) // The detected format.

	return m, nil // Placeholder
}
