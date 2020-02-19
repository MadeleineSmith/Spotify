package handlers

import (
	"Spotify/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type CreatePlaylistHandler struct {
	HTTPClient *http.Client
}

type CurrentUserResponse struct {
	ID string `json:"id"`
}

type CreatePlaylistRequest struct {
	Date string `json:"date"`
}

func (h CreatePlaylistHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	spotifyAccessToken := req.Header.Get("Authorization")
	userID := h.getUserID(spotifyAccessToken)

	inputBodyBytes, _ := ioutil.ReadAll(req.Body)
	var createPlaylistRequest CreatePlaylistRequest
	json.Unmarshal(inputBodyBytes, &createPlaylistRequest)

	// might be a good idea to write a custom unmarshal function for the time type but cba
	layout := "20060102"
	incomingTime, _ := time.Parse(layout, createPlaylistRequest.Date)

	layoutUS := "January 2 2006"
	playlistName := fmt.Sprintf("%s Chart", incomingTime.Format(layoutUS))

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
