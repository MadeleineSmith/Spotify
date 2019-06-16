package models

type Track struct {
	URL string `json:"uri"`
	TrackName string `json:"trackName"`
	ArtistName string `json:"artistName"`
}