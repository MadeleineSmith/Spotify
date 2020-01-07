package models

type Chart struct {
	Date   string  `json:"date"`
	Tracks []Track `json:"tracks"`
}
