package handlers

import (
	. "Spotify/constants"
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
	Tracks Track `json:"tracks"`
}

type Track struct {
	Items []TrackItem `json:"items"`
}

type TrackItem struct {
	URL string `json:"uri"`
	TrackName string `json:"trackName"`
	ArtistName string `json:"artistName"`
}

func (h SearchHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	inputBodyBytes, _ := ioutil.ReadAll(req.Body)
	trackItem := new(TrackItem)
	json.Unmarshal(inputBodyBytes, &trackItem)

	spotifyRequestURL := buildURL(trackItem)

	spotifyRequest, _ := http.NewRequest(http.MethodGet, spotifyRequestURL, nil)
	spotifyRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AUTHORIZATION_TOKEN))

	spotifyResponse, _ := h.HTTPClient.Do(spotifyRequest)

	spotifyBodyBytes, _ := ioutil.ReadAll(spotifyResponse.Body)
	spotifyResponseBody := new(SpotifyResponse)
	json.Unmarshal(spotifyBodyBytes, &spotifyResponseBody)

	trackItem.URL = spotifyResponseBody.Tracks.Items[0].URL

	data, _ := json.Marshal(trackItem)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func buildURL(trackItem *TrackItem) string {
	baseUrl, _ := url.Parse("https://api.spotify.com/v1/search")
	params := url.Values{}
	params.Add("q", fmt.Sprintf("track:%s", trackItem.TrackName))
	params.Add("q", fmt.Sprintf("artist:%s", trackItem.ArtistName))
	params.Add("limit", "1")
	params.Add("type", "track")
	baseUrl.RawQuery = params.Encode()

	return baseUrl.String()
}