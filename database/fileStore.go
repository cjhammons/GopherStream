package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/dhowden/tag"
)

const STORAGE_DIR = "storage"
const ALBUM_ART_DIR = "album-art"
const ARTIST_ART_DIR = "artist-art"

func SaveAlbumArt(picture *tag.Picture, albumID int64) (string, error) {
	// Create a filename for the image
	storageDir := filepath.Join(STORAGE_DIR, ALBUM_ART_DIR)
	fileName := fmt.Sprintf("album-art-%d.%s", albumID, picture.Ext)
	filePath := filepath.Join(storageDir, fileName)

	// Check if the storage directory exists, if not, create it
	if _, err := os.Stat(storageDir); os.IsNotExist(err) {
		if err := os.MkdirAll(storageDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create storage directory: %v", err)
		}
	} else if err != nil {
		// An error occurred when checking the directory
		return "", fmt.Errorf("error checking storage directory: %v", err)
	}

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		log.Printf("Album art already exists: %v", filePath)
		return filePath, nil
	} else if !os.IsNotExist(err) {
		// If the error is not that the file does not exist, return the error
		return "", err
	}

	// Save the image data to a file
	err := os.WriteFile(filePath, picture.Data, 0644)
	if err != nil {
		return "", err
	}
	log.Printf("Saved album art to: %v", filePath)
	return filePath, nil
}
