package database

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
	artistID, err := InsertOrUpdateArtist(db, metadata)
	if err != nil {
		return err
	}

	albumID, err := InsertOrUpdateAlbum(db, metadata, artistID)
	if err != nil {
		return err
	}

	err = InsertOrUpdateAlbumArt(db, metadata, albumID)
	if err != nil {
		return err
	}

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
