package handlers

import (
	"Spotify/models"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
	"net/http"
)

type ScrapeChartsHandler struct {
	HTTPClient *http.Client
}

func (h ScrapeChartsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	chartDate := vars["chart_date"]

	tracks := []models.Track{}

	c := colly.NewCollector()

	c.OnHTML("table.chart-positions", func(table *colly.HTMLElement) {
		table.ForEach("tr", func(_ int, row *colly.HTMLElement) {
			title := row.ChildText("div.title")
			artist := row.ChildText("div.artist")

			// feels slightly hacky... but hey
			if title != "" && artist != "" {
				tracks = append(tracks, models.Track{
					TrackName: title,
					ArtistName: artist,
				})
			}
		})
	})

	c.Visit(fmt.Sprintf("https://www.officialcharts.com/charts/singles-chart/%s/7501/", chartDate))

	data, _ := json.Marshal(tracks)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}