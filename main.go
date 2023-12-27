package main

import (
	"cjhammons.com/goaudio/config"
	"fmt"
	"github.com/dhowden/tag"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	cfg.Print()

	http.HandleFunc("/stream", streamHandler)
	log.Println("Server is running on http://localhost:8080/stream")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// streamHandler handles the HTTP request to stream the audio file
func streamHandler(w http.ResponseWriter, r *http.Request) {
	// Open the source audio file
	sourceFile, err := os.Open("synthwave-loop.mp3")
	if err != nil {
		log.Printf("Error opening source file: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer sourceFile.Close()

	// Set the appropriate header to inform the client that the content is an audio file
	w.Header().Set("Content-Type", "audio/mpeg")

	m, err := tag.ReadFrom(sourceFile)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(m.Format()) // The detected format.
	fmt.Print(m.Title())  // The title of the track (see Metadata interface for more details).
	// Stream the audio from source to the HTTP response
	_, err = io.Copy(w, sourceFile)
	if err != nil {
		log.Printf("Error streaming file: %v", err)
	}
}
