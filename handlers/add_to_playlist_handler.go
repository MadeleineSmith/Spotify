package handlers

import (
	"Spotify/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type AddToPlaylistHandler struct {
	HTTPClient *http.Client
}

type SpotifyAddToPlaylistRequest struct {
	URIs []string `json:"uris"`
}

func (h AddToPlaylistHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	spotifyAccessToken := req.Header.Get("Authorization")

	vars := mux.Vars(req)
	playlistID := vars["playlist_id"]

	inputBodyBytes, _ := ioutil.ReadAll(req.Body)
	var inputTrackData []*models.Track
	json.Unmarshal(inputBodyBytes, &inputTrackData)

	var spotifyURIs []string
	for _, track := range inputTrackData {
		if track.URI != "" {
			spotifyURIs = append(spotifyURIs, track.URI)
		}
	}

	spotifyTrackRequest := SpotifyAddToPlaylistRequest{
		URIs: spotifyURIs,
	}

	data, _ := json.Marshal(spotifyTrackRequest)

	spotifyRequestURL := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", playlistID)

	spotifyRequest, _ := http.NewRequest(http.MethodPost, spotifyRequestURL, bytes.NewBuffer(data))

	spotifyRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", spotifyAccessToken))
	h.HTTPClient.Do(spotifyRequest)
}
