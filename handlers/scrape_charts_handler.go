package handlers

import (
	"Spotify/models"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
)

type ScrapeChartsHandler struct {
	HTTPClient *http.Client
}

func (h ScrapeChartsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	year := vars["year"]

	randomDateString := getRandomDateString(year)

	chart := models.Chart{}
	chart.Date = randomDateString

	chart.Tracks = []models.Track{}

	c := colly.NewCollector()

	c.OnHTML("table.chart-positions", func(table *colly.HTMLElement) {
		table.ForEach("tr", func(_ int, row *colly.HTMLElement) {
			title := row.ChildText("div.title")
			artist := row.ChildText("div.artist")

			// feels slightly hacky... but hey
			if title != "" && artist != "" {
				chart.Tracks = append(chart.Tracks, models.Track{
					TrackName:  title,
					ArtistName: artist,
				})
			}
		})
	})

	officialChartsURL := fmt.Sprintf("https://www.officialcharts.com/charts/singles-chart/%s/7501/", randomDateString)

	c.Visit(officialChartsURL)

	data, _ := json.Marshal(chart)

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
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
