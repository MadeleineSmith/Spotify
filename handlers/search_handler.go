package handlers

import (
	"Spotify/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
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
	spotifyAccessToken := req.Header.Get("Authorization")

	inputBodyBytes, _ := ioutil.ReadAll(req.Body)

	var chart models.Chart

	//var inputTrackData []*models.Track
	json.Unmarshal(inputBodyBytes, &chart)

	for _, track := range chart.Tracks {
		spotifyRequestURL := h.buildURL(track)
		h.makeSpotifySearchRequest(spotifyRequestURL, track, spotifyAccessToken)
	}

	// TODO - for some reason not adding the URI
	data, _ := json.Marshal(chart)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func (h SearchHandler) buildURL(track *models.Track) string {
	lowerCaseArtist := strings.ToLower(track.ArtistName)
	lowerCaseTrackName := strings.ToLower(track.TrackName)

	artistReplacer := strings.NewReplacer(
		" ft ", " ",
		"/", " ",
		"&", " ",
		"chipmunk", "chip", // ok
		"will i am", "will.i.am", // ok
		"lily rose cooper", "lily allen", // ok
	)

	// Official Charts corrections
	trackReplacer := strings.NewReplacer(
		"you got the love", "you've got the love",
		"sos (let the music play)", "s.o.s. (let the music play)",
		"she's got me dancin", "she's got me dancing",
		"just the way you are (amazing)", "just the way you are", // ok
	)

	artistWithReplacements := artistReplacer.Replace(lowerCaseArtist)
	trackWithReplacements := trackReplacer.Replace(lowerCaseTrackName)

	// doing two separate regular expressions as golang does not allow negative lookaheads
	pinkFloydRe := regexp.MustCompile(`pink floyd`)
	pinkRe := regexp.MustCompile(`pink`)

	if pinkRe.MatchString(artistWithReplacements) && !pinkFloydRe.MatchString(artistWithReplacements) {
		artistWithReplacements = pinkRe.ReplaceAllString(artistWithReplacements, `p!nk`)
	}

	re := regexp.MustCompile(`^et$`)
	trackWithReplacements = re.ReplaceAllString(trackWithReplacements, `e.t.`)

	baseURL, _ := url.Parse("https://api.spotify.com/v1/search")
	params := url.Values{}

	massiveString := fmt.Sprintf("%s artist:%s", trackWithReplacements, artistWithReplacements)
	params.Add("q", massiveString)

	params.Add("limit", "1")
	params.Add("type", "track")
	baseURL.RawQuery = params.Encode()

	return baseURL.String()
}

func (h SearchHandler) makeSpotifySearchRequest(url string, track *models.Track, accessToken string) {
	spotifyRequest, _ := http.NewRequest(http.MethodGet, url, nil)
	spotifyRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	spotifyResponse, _ := h.HTTPClient.Do(spotifyRequest)

	spotifyBodyBytes, _ := ioutil.ReadAll(spotifyResponse.Body)
	spotifyResponseBody := new(SpotifyResponse)
	json.Unmarshal(spotifyBodyBytes, &spotifyResponseBody)

	// TODO - do I want to remove this track from the JSON I return?
	if len(spotifyResponseBody.Tracks.Items) == 1 {
		track.URI = spotifyResponseBody.Tracks.Items[0].URI
	}
}
