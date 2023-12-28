package models

type Album struct {
	ID          int64
	Title       string
	Artist      *Artist
	ReleastDate string
	ArtFilePath string
}
