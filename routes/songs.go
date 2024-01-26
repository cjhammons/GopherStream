package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"cjhammons.com/gopher-stream/database"
)

func GetSongHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		songHandler(db, w, r)
	}
}

func songHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	songs, err := database.GetAllSongs(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(songs)

}
