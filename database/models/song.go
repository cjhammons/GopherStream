package models

// type Song struct {
// 	ID         int64
// 	Title      string
// 	ArtistID   int64
// 	AlbumID    int64
// 	GenreID    int64
// 	TrackNum   int64
// 	FilePath   string
// 	FileFormat string
// }

type Song struct {
	ID         int64
	Title      string
	Artist     string
	Album      string
	Genre      string
	TrackNum   int64
	FilePath   string
	FileFormat string
}
