package main

import (
	"Spotify/handlers"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
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

	port := os.Getenv("PORT")

	fmt.Printf("Running on port %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), CORS(router)))
}

// TODO - hmmmm, for some reason I cannot get rs/cors to work so am using this instead:
func CORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// allowing all origins
		// TODO - mebs limit this down in the future
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		h.ServeHTTP(w, r)
	})
}
