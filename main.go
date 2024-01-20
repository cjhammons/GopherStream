package main

import (
	"log"

	"cjhammons.com/goaudio/config"
	"cjhammons.com/goaudio/database"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	cfg.Print()

	db, err := database.InitializeDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	database.ScanLibrary(cfg.LibraryDirectory, db)

	// http.HandleFunc("/stream", streamHandler)
	// log.Println("Server is running on http://localhost:8080/stream")
	// log.Fatal(http.ListenAndServe(":8080", nil))
}
