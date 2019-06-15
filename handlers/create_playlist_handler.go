package handlers

import (
	. "Spotify/constants"
	"bytes"
	"fmt"
	"net/http"
)

type CreatePlaylistHandler struct{
	HTTPClient *http.Client
}

func (h CreatePlaylistHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {}

	playlistName := req.Form.Get("name")

	jsonString := fmt.Sprintf(
		`{
  "name": "%s",
  "public": false
}`, playlistName)

	request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", USER_ID), bytes.NewBuffer([]byte(jsonString)))
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AUTHORIZATION_TOKEN))

	h.HTTPClient.Do(request)
}