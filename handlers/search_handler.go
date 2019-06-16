package handlers

import (
	. "Spotify/constants"
	"Spotify/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type SearchHandler struct {
	HTTPClient *http.Client
}

type SpotifyResponse struct {
	Tracks TrackItems `json:"tracks"`
}

type TrackItems struct {
	Items []models.Track `json:"items"`
}

func (h SearchHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	inputBodyBytes, _ := ioutil.ReadAll(req.Body)
	track := new(models.Track)
	json.Unmarshal(inputBodyBytes, &track)

	spotifyRequestURL := buildURL(track)

	spotifyRequest, _ := http.NewRequest(http.MethodGet, spotifyRequestURL, nil)
	spotifyRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AUTHORIZATION_TOKEN))

	spotifyResponse, _ := h.HTTPClient.Do(spotifyRequest)

	spotifyBodyBytes, _ := ioutil.ReadAll(spotifyResponse.Body)
	spotifyResponseBody := new(SpotifyResponse)
	json.Unmarshal(spotifyBodyBytes, &spotifyResponseBody)

	track.URL = spotifyResponseBody.Tracks.Items[0].URL

	data, _ := json.Marshal(track)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func buildURL(track *models.Track) string {
	baseUrl, _ := url.Parse("https://api.spotify.com/v1/search")
	params := url.Values{}
	params.Add("q", fmt.Sprintf("track:%s", track.TrackName))
	params.Add("q", fmt.Sprintf("artist:%s", track.ArtistName))
	params.Add("limit", "1")
	params.Add("type", "track")
	baseUrl.RawQuery = params.Encode()

	return baseUrl.String()
}