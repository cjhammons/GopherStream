package models

type Album struct {
	ID          int64
	Title       string
	ArtistID    int64 // Foreign key for the Artist
	ReleaseDate string
	ArtFilePath string
}
