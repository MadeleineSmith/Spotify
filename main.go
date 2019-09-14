package main

import (
	. "Spotify/constants"
	"Spotify/handlers"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"strconv"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	client := &http.Client{}

	loginUserHandler := handlers.LoginUserHandler{
		HTTPClient: client,
	}
	callbackHandler := handlers.CallbackHandler{
		HTTPClient: client,
	}

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

	router.NewRoute().Path("/login").Handler(loginUserHandler)
	router.NewRoute().Path("/callback").Handler(callbackHandler)

	router.NewRoute().Path("/user/playlists").Handler(createPlaylistHandler)
	router.NewRoute().Path("/charts/{year}").Handler(scrapeChartsHandler)
	router.NewRoute().Path("/search").Handler(searchHandler)
	router.NewRoute().Path("/playlists/{playlist_id}/tracks").Handler(addToPlaylistHandler)

	httpHandler := cors.Default().Handler(router)

	fmt.Printf("Running on port: %s\n", strconv.Itoa(PORT))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", strconv.Itoa(PORT)), httpHandler))
}
