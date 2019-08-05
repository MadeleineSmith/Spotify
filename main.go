package main

import (
	"Spotify/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	client := &http.Client{}

	scrapeChartsHandler := handlers.ScrapeChartsHandler{}

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
	router.NewRoute().Path("/charts/{year}").Handler(scrapeChartsHandler)
	router.NewRoute().Path("/search").Handler(searchHandler)
	router.NewRoute().Path("/playlists/{playlist_id}/tracks").Handler(addToPlaylistHandler)

	// TODO - should probs research this CORS stuff later
	httpHandler := cors.Default().Handler(router)

	log.Fatal(http.ListenAndServe(":6584", httpHandler))
}
