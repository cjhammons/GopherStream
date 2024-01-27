package routes

import (
	"database/sql"
	"net/http"
	"text/template"

	"cjhammons.com/gopher-stream/database"
)

func GetSongHtmxHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		songHandler(db, w, r)
	}
}

func songHtmxHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	songs, err := database.GetAllSongs(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	tmpl := template.Must(template.ParseFiles("templates/songs.html"))
	tmpl.Execute(w, songs)
}
