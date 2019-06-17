package main

import (
	"Spotify/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	client := &http.Client{}

	// TODO - check with Tom about base class with `client` property?
	createPlaylistHandler := handlers.CreatePlaylistHandler{
		HTTPClient: client,
	}
	searchHandler := handlers.SearchHandler{
		HTTPClient: client,
	}
	addToPlaylistHandler := handlers.AddToPlaylistHandler{
		HTTPClient: client,
	}

	router.NewRoute().Path("/users/{user_id}/playlists").Handler(createPlaylistHandler)
	router.NewRoute().Path("/search").Handler(searchHandler)
	router.NewRoute().Path("/playlists/{playlist_id}/tracks").Handler(addToPlaylistHandler)

	log.Fatal(http.ListenAndServe(":6584", router))
}
