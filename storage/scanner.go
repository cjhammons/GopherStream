package storage

import (
	"cjhammons.com/goaudio/models"
	"database/sql"
	"dhawden/tag"
	"log"
	"os"
	"path/filepath"
)

var allowedFileFormats = []string{".mp3", ".flac"}

func ScanLibrary(directory string, db *sql.DB) error {
	log.Println("Beginning library svan...")
	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && isAllowedAudioFile(path) {
			err := processAudioFile(path, db)
			if err != nil {
				return err
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
		log.Fatalf("Error extracting metadata: %v", err)
		return err
	}
	// Insert metadata into database
	err = InsertArtist(db, metadata.Artist(), "") //todo filepath is blank
	if err != nil {
		return err
	}
	// Similarly, insert Album, Genre, and Song

	return nil
}

/*
Note Metadata interface from dhawden/tag:

	type Metadata interface {
		Format() Format
		FileType() FileType

		Title() string
		Album() string
		Artist() string
		AlbumArtist() string
		Composer() string
		Genre() string
		Year() int

		Track() (int, int) // Number, Total
		Disc() (int, int) // Number, Total

		Picture() *Picture // Artwork
		Lyrics() string
		Comment() string

		Raw() map[string]interface{} // NB: raw tag names are not consistent across formats.
	}
*/
func extractMetadata(filePath string) (*tag.Metadata, error) {
	m, err := tag.ReadFrom(filePath)
	if err != nil {
		return nil, err
	}
	log.Print(m.Format()) // The detected format.

	return m, nil // Placeholder
}
