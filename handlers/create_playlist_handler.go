package handlers

import (
	"Spotify/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
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
	MinimumYear  *int        `json:"minYear"`
	SpecificDate *customDate `json:"date"`
}

type customDate struct {
	time.Time
}

func (sd *customDate) UnmarshalJSON(input []byte) error {
	strInput := string(input)
	strInput = strings.Trim(strInput, `"`)
	newTime, err := time.Parse("2006-01-02", strInput)
	if err != nil {
		return err
	}

	sd.Time = newTime
	return nil
}

// need to verify on the dates coming in that not before or after Official Charts range
func (h CreatePlaylistHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		h.handlePost(w, req)
	}
}

func (h CreatePlaylistHandler) handlePost(w http.ResponseWriter, req *http.Request) {
	spotifyAccessToken := req.Header.Get("Authorization")
	userID := h.getUserID(spotifyAccessToken)

	inputBodyBytes, _ := ioutil.ReadAll(req.Body)
	var createPlaylistRequest CreatePlaylistRequest
	json.Unmarshal(inputBodyBytes, &createPlaylistRequest)

	var playlistName string
	var playlistDate time.Time

	// if only year provided
	if createPlaylistRequest.MinimumYear != nil && createPlaylistRequest.SpecificDate == nil {
		playlistDate = getRandomDate(*createPlaylistRequest.MinimumYear)
	} else if createPlaylistRequest.SpecificDate != nil && createPlaylistRequest.MinimumYear == nil {
		playlistDate = createPlaylistRequest.SpecificDate.Time
	}

	layoutUK := "2 January 2006"
	playlistName = fmt.Sprintf("%s Chart", playlistDate.Format(layoutUK))

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
	playlist.Date = playlistDate.Format("2006-01-02")

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

// should probs unit test...
func getRandomDate(minimumYear int) time.Time {
	var min int64
	var max int64

	currentDate := time.Now()
	currentYear := currentDate.Year()

	if minimumYear == 1952 {
		// Records start on 14/11/1952
		min = time.Date(minimumYear, 11, 14, 0, 0, 0, 0, time.UTC).Unix()
	} else {
		min = time.Date(minimumYear, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
	}

	if minimumYear == currentYear {
		max = time.Date(currentYear, currentDate.Month(), currentDate.Day(), 23, 59, 59, 999999999, time.UTC).Unix()
	} else {
		max = time.Date(currentYear, 12, 31, 23, 59, 59, 999999999, time.UTC).Unix()
	}

	secondsBetweenDates := max - min

	seed := rand.NewSource(time.Now().UnixNano())
	seededRand := rand.New(seed)

	randomDate := min + seededRand.Int63n(secondsBetweenDates)

	// using UTC to prevent overlapping to subsequent days due to differing time zones
	// found to be necessary when year == currentDate.Year()
	return time.Unix(randomDate, 0).UTC()
}
