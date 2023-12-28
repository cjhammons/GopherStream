package models

type Song struct {
	ID       int64
	Title    string
	Artist   *Artist
	Album    *Album
	Genre    *Genre
	FilePath string
}
