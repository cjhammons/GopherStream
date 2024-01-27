package main

import (
	"log"
	"os"
	"time"

	"net/http"

	"cjhammons.com/gopher-stream/config"
	"cjhammons.com/gopher-stream/database"
	"cjhammons.com/gopher-stream/routes"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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

	router := mux.NewRouter()
	router.HandleFunc("/songs", routes.GetSongHtmxHandler(db)).Methods("GET")

	// Wrap router with Gorilla Handlers for additional functionality like Logging
	loggingRouter := handlers.LoggingHandler(os.Stdout, router)

	// Start server
	srv := &http.Server{
		Handler:      loggingRouter,
		Addr:         "0.0.0.0:6969",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Starting server on port 6969")

	log.Fatal(srv.ListenAndServe())

}
