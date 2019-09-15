package handlers

import (
	"Spotify/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CreatePlaylistHandler struct {
	HTTPClient *http.Client
}

type CurrentUserResponse struct {
	ID string `json:"id"`
}

func (h CreatePlaylistHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	spotifyAccessToken := req.Header.Get("Authorization")
	userID := h.getUserID(spotifyAccessToken)

	req.ParseForm()
	// TODO - error if no name is provided?
	playlistName := req.Form.Get("name")

	// TODO - following feels a bit lazy
	jsonString := fmt.Sprintf(
		`{
  "name": "%s",
  "public": false
}`, playlistName)

	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", userID), bytes.NewBuffer([]byte(jsonString)))
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", spotifyAccessToken))
	spotifyResponse, _ := h.HTTPClient.Do(request)

	spotifyBodyBytes, _ := ioutil.ReadAll(spotifyResponse.Body)

	playlist := new(models.Playlist)
	json.Unmarshal(spotifyBodyBytes, &playlist)

	data, _ := json.Marshal(playlist)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h CreatePlaylistHandler) getUserID(spotifyAccessToken string) string {
	request, _ := http.NewRequest(http.MethodGet, "https://api.spotify.com/v1/me", nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", spotifyAccessToken))

	spotifyResponse, _ := h.HTTPClient.Do(request)

	spotifyResponseBodyBytes, _ := ioutil.ReadAll(spotifyResponse.Body)

	var currentUser CurrentUserResponse
	json.Unmarshal(spotifyResponseBodyBytes, &currentUser)

	return currentUser.ID
}
