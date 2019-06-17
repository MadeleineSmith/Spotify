package models

type Track struct {
	URI        string `json:"uri"`
	TrackName  string `json:"trackName"`
	ArtistName string `json:"artistName"`
}