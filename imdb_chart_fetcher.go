package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var scrapingCount int

type movies struct {
	Title            string `json:"title"`
	MovieReleaseYear string `json:"movie_release_year"`
	ImdbRating       string `json:"imdb_rating"`
	Summary          string `json:"summary"`
	Duration         string `json:"duration"`
	Genre            string `json:"genre"`
}

var movielist []movies

func main() {
	var imdbURL string

	imdbURL = os.Args[1]
	scrapingCount, _ = strconv.Atoi(os.Args[2])

	// Make HTTP GET request
	response, err := http.Get(imdbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	document.Find(".lister .chart .lister-list tr .titleColumn a").EachWithBreak(processElement)

	// Result
	resultset, _ := json.Marshal(movielist)
	fmt.Printf(string(resultset))
}

// Action for each movies in the top moveis list
func processElement(index int, element *goquery.Selection) bool {
	href, exists := element.Attr("href")
	if exists {
		if scrapingCount > 0 {
			scrapingCount--
			getmovieData(href)
		} else {
			return false
		}
	}
	return true
}

// Get data from each individual movie page
func getmovieData(url string) {
	// Make HTTP GET request
	response, err := http.Get("https://www.imdb.com/" + url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	// Movie Data

	// Movie Name
	movieTitle, _ := document.Find("#star-rating-widget").Attr("data-title")
	// Movie Year
	movieYear := document.Find("#titleYear a").Text()
	// IMDB Rating
	movieIMDBRating := document.Find(".imdbRating .ratingValue strong span").Text()
	// Movie Summary
	movieSummary := strings.TrimSpace(document.Find(".plot_summary .summary_text").Text())
	// Movie Time
	movieDuration := strings.TrimSpace(document.Find(".title_wrapper .subtext time").Text())
	// Movie Genres
	movieInfo := strings.Split(document.Find(".title_wrapper .subtext").Text(), "|")
	movieGenres := strings.Split(movieInfo[2], ",")
	for i, s := range movieGenres {
		movieGenres[i] = strings.TrimSpace(s)
	}

	// Save each movie data
	movieData := movies{movieTitle, movieYear, movieIMDBRating, movieSummary, movieDuration, strings.Join(movieGenres, ", ")}
	movielist = append(movielist, movieData)
}
