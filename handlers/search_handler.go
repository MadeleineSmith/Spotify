package handlers

import (
	. "Spotify/constants"
	"Spotify/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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

	var inputTrackData []*models.Track
	json.Unmarshal(inputBodyBytes, &inputTrackData)

	for _, track := range inputTrackData {
		spotifyRequestURL := h.buildURL(track)
		h.makeSpotifySearchRequest(spotifyRequestURL, track)
	}

	data, _ := json.Marshal(inputTrackData)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h SearchHandler) buildURL(track *models.Track) string {
	artistWithoutFt := strings.Replace(track.ArtistName, " FT ", " ", -10)

	baseURL, _ := url.Parse("https://api.spotify.com/v1/search")
	params := url.Values{}

	query := strings.ToLower(fmt.Sprintf("%s %s", track.TrackName, artistWithoutFt))

	params.Add("q", query)

	params.Add("limit", "1")
	params.Add("type", "track")
	baseURL.RawQuery = params.Encode()

	return baseURL.String()
}

func (h SearchHandler) makeSpotifySearchRequest(url string, track *models.Track) {
	spotifyRequest, _ := http.NewRequest(http.MethodGet, url, nil)
	spotifyRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", AUTHORIZATION_TOKEN))

	spotifyResponse, _ := h.HTTPClient.Do(spotifyRequest)

	spotifyBodyBytes, _ := ioutil.ReadAll(spotifyResponse.Body)
	spotifyResponseBody := new(SpotifyResponse)
	json.Unmarshal(spotifyBodyBytes, &spotifyResponseBody)

	// TODO - do I want to remove this track from the JSON I return?
	if len(spotifyResponseBody.Tracks.Items) == 1 {
		track.URI = spotifyResponseBody.Tracks.Items[0].URI
	}
}