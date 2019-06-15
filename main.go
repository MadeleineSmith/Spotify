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

	// TODO - change this to be more restful i.e POST /playlist
	router.NewRoute().Path("/create-playlist").Handler(createPlaylistHandler)
	router.NewRoute().Path("/search").Handler(searchHandler)

	log.Fatal(http.ListenAndServe(":6584", router))
}
