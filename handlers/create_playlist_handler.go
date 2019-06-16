package handlers

import (
	. "Spotify/constants"
	"Spotify/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type CreatePlaylistHandler struct{
	HTTPClient *http.Client
}

func (h CreatePlaylistHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userID := vars["user_id"]

	if err := req.ParseForm(); err != nil {}

	playlistName := req.Form.Get("name")

	jsonString := fmt.Sprintf(
		`{
  "name": "%s",
  "public": false
}`, playlistName)

	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", userID), bytes.NewBuffer([]byte(jsonString)))
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AUTHORIZATION_TOKEN))
	spotifyResponse, _ := h.HTTPClient.Do(request)

	spotifyBodyBytes, _ := ioutil.ReadAll(spotifyResponse.Body)

	playlist := new(models.Playlist)
	json.Unmarshal(spotifyBodyBytes, &playlist)

	data, _ := json.Marshal(playlist)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}