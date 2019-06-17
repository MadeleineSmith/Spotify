package handlers

import (
	. "Spotify/constants"
	"Spotify/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/url"
)

type AddToPlaylistHandler struct {
	HTTPClient *http.Client
}

func (h AddToPlaylistHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	playlistID := vars["playlist_id"]

	inputBodyBytes, _ := ioutil.ReadAll(req.Body)
	track := new(models.Track)
	json.Unmarshal(inputBodyBytes, &track)

	spotifyRequestURL := h.buildURL(playlistID, track)
	spotifyRequest, _ := http.NewRequest(http.MethodPost, spotifyRequestURL, nil)
	spotifyRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AUTHORIZATION_TOKEN))

	h.HTTPClient.Do(spotifyRequest)
}

// TODO - consider changing to JSON body for track URIs instead of query parameters
func (h AddToPlaylistHandler) buildURL(playlistID string, track *models.Track) string {
	baseURL, _ := url.Parse(fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistID))
	params := url.Values{}
	params.Add("uris", track.URI)

	baseURL.RawQuery = params.Encode()

	return baseURL.String()
}
