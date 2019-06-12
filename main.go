package main

import (
	"Spotify/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	createPlaylistHandler := handlers.CreatePlaylistHandler{}

	router.NewRoute().Path("/create-playlist").Handler(createPlaylistHandler)

	log.Fatal(http.ListenAndServe(":6584", router))
}
