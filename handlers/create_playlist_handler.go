package handlers

import (
	"Spotify/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type CreatePlaylistHandler struct {
	HTTPClient *http.Client
}

type CurrentUserResponse struct {
	ID string `json:"id"`
}

// Add *name* in future
type CreatePlaylistRequest struct {
	Year string `json:"year"`
}

func (h CreatePlaylistHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	spotifyAccessToken := req.Header.Get("Authorization")
	userID := h.getUserID(spotifyAccessToken)

	inputBodyBytes, _ := ioutil.ReadAll(req.Body)
	var createPlaylistRequest CreatePlaylistRequest
	json.Unmarshal(inputBodyBytes, &createPlaylistRequest)

	randomDateString := getRandomDateString(createPlaylistRequest.Year)

	layout := "20060102"
	randomDate, _ := time.Parse(layout, randomDateString)

	layoutUS := "January 2 2006"
	playlistName := fmt.Sprintf("%s Chart", randomDate.Format(layoutUS))

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
	playlist.Date = randomDateString

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

// TODO - should probs test these two functions
func getRandomDateString(yearString string) string {
	year, _ := strconv.Atoi(yearString)

	randomDateInYear := getRandomDateInYear(year)

	paddedMonth := fmt.Sprintf("%02d", randomDateInYear.Month())
	paddedDay := fmt.Sprintf("%02d", randomDateInYear.Day())
	randomDateString := fmt.Sprintf("%s%s%s", yearString, paddedMonth, paddedDay)

	return randomDateString
}

func getRandomDateInYear(year int) time.Time {
	var min int64
	var max int64

	currentDate := time.Now()

	if year == 1952 {
		// Records start on 14/11/1952
		min = time.Date(year, 11, 14, 0, 0, 0, 0, time.UTC).Unix()
	} else {
		min = time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	}

	if year == currentDate.Year() {
		max = time.Date(year, currentDate.Month(), currentDate.Day(), 23, 59, 59, 999999999, time.UTC).Unix()
	} else {
		max = time.Date(year, 12, 31, 23, 59, 59, 999999999, time.UTC).Unix()
	}

	secondsBetweenDates := max - min

	seed := rand.NewSource(time.Now().UnixNano())
	seededRand := rand.New(seed)

	randomDate := min + seededRand.Int63n(secondsBetweenDates)

	// using UTC to prevent overlapping to subsequent days due to differing time zones
	// found to be necessary when year == currentDate.Year()
	return time.Unix(randomDate, 0).UTC()
}
