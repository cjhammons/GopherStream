package storage

import (
	"database/sql"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/dhowden/tag"
)

var allowedFileFormats = []string{".mp3", ".flac"}

func ScanLibrary(directory string, db *sql.DB) error {
	// If the filepath contains '~' then Expand the '~' to the user's home directory
	// This is a workaround for the fact that the filepath package does not expand '~'
	if directory[:2] == "~/" {
		usr, err := user.Current()
		if err != nil {
			log.Fatalf("Cannot get current user: %v", err)
		}
		directory = filepath.Join(usr.HomeDir, directory[2:])
	}
	log.Println("Beginning library scan...")
	log.Println("Scanning directory: ", directory)
	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing file %s: %v", path, err)
			return err
		}
		if !info.IsDir() && isAllowedAudioFile(path) {
			if err := processAudioFile(path, db); err != nil {
				log.Printf("Error processing audio file %s: %v", path, err)
			} else {
				log.Printf("File Successfuly Processed: %v", path)
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

	_, err = InsertOrUpdateSong(db, metadata, albumID, artistID, genreID, filePath)
	if err != nil {
		return err
	}

	return nil
}

func InsertOrUpdateArtist(db *sql.DB, metadata tag.Metadata, filePath string) (int64, error) {
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

func InsertOrUpdateAlbum(db *sql.DB, metadata tag.Metadata, artistID int64, filePath string) (int64, error) {
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
                release_date = ?, 
                art_file_path = ? 
            WHERE 
                id = ?`
		_, err := db.Exec(updateQuery, metadata.Year(), filePath, albumID)
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
	`
	err := db.QueryRow(existsQuery, metadata.Title(), albumID).Scan(&songID)
	if err == nil {
		log.Printf("Song \"%v\" already exists. Updating...", metadata.Title())

		// Update existing song
		updateQuery := `
            UPDATE songs 
            SET 
                artist_id = ?, 
                genre_id = ?, 
                track_number = ?, 
                file_path = ?, 
                file_format = ? 
            WHERE 
                id = ?`
		_, err := db.Exec(updateQuery, artistID, genreID, trackNumber, filePath, metadata.FileType(), songID)
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

func extractMetadata(filePath string) (tag.Metadata, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	m, err := tag.ReadFrom(file)
	if err != nil {
		return nil, err
	}

	return m, nil
}
